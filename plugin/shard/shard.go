package main

import (
	"os"

	"github.com/discentem/plugexperiments/commons"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

// Here is a real implementation of Greeter
type ShardEx struct {
	logger hclog.Logger
}

func (s *ShardEx) Get() (string, error) {
	s.logger.Debug("message from GreeterHello.Greet")
	return "shard!", nil
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

	sharder := &ShardEx{
		logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"shard": &commons.ShardPlugin{Impl: sharder},
	}

	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
