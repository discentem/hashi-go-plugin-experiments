package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/discentem/plugexperiments/commons"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var (
	handshakeConfig = plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "BASIC_PLUGIN",
		MagicCookieValue: "hello",
	}

	pluginMap = map[string]plugin.Plugin{
		"greeter":    &commons.GreeterPlugin{},
		"greeterToo": &commons.GreeterPlugin{},
	}
)

func main() {
	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	plugs, err := plugin.Discover("*", "./plugin/")
	if err != nil {
		log.Fatal(err)
	}
	for i := range plugs {
		if strings.HasSuffix(plugs[i], ".go") {
			continue
		}
		// We're a host! Start by launching the plugin process.
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: handshakeConfig,
			Plugins:         pluginMap,
			Cmd:             exec.Command(plugs[i]),
			Logger:          logger,
		})
		defer client.Kill()

		// Connect via RPC
		rpcClient, err := client.Client()
		if err != nil {
			log.Fatal(err)
		}

		// Request the plugin

		for _, p := range []string{"greeter", "greeterToo"} {
			raw, err := rpcClient.Dispense(p)
			if err != nil {
				log.Fatal(err)
			}

			// We should have a Greeter now! This feels like a normal interface
			// implementation but is in fact over an RPC connection.
			greeter := raw.(commons.Greeter)
			s, err := greeter.Greet()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(s)
			f, err := greeter.GreetFancy()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(f)
		}
	}

}
