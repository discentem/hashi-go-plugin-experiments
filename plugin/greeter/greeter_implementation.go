package main

import (
	"os"

	"github.com/discentem/plugexperiments/commons"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

// Here is a real implementation of Greeter
type GreeterHello struct {
	logger hclog.Logger
}

func (g *GreeterHello) Greet() (string, error) {
	g.logger.Debug("message from GreeterHello.Greet")
	return "Hello!", nil
}

func (g *GreeterHello) GreetFancy() (string, error) {
	g.logger.Debug("message from GreeterHello.GreetFancy")
	return "Hello, fancy pants!", nil
}

type GreeterHelloToo struct {
	logger hclog.Logger
}

func (g *GreeterHelloToo) Greet() (string, error) {
	g.logger.Debug("message from GreeterHelloToo.Greet")
	return "Hello!", nil
}

func (g *GreeterHelloToo) GreetFancy() (string, error) {
	g.logger.Debug("message from GreeterHelloToo.GreetFancy")
	return "Hello, fancy pants!", nil
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	greeter := &GreeterHello{
		logger: logger,
	}
	greeterToo := &GreeterHelloToo{
		logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"greeter":    &commons.GreeterPlugin{Impl: greeter},
		"greeterToo": &commons.GreeterPlugin{Impl: greeterToo},
	}

	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
