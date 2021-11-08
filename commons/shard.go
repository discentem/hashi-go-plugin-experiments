package commons

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Greeter is the interface that we're exposing as a plugin.
type Shard interface {
	Get() (string, error)
}

// Here is an implementation that talks over RPC
type ShardRPC struct{ client *rpc.Client }

func (s *ShardRPC) Get() (string, error) {
	var resp string
	err := s.client.Call("Plugin.Get", new(interface{}), &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		return "", err
	}

	return resp, nil
}

// Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type ShardRPCServer struct {
	// This is the real implementation
	Impl Shard
}

func (s *ShardRPCServer) Get(args interface{}, resp *string) error {
	*resp, _ = s.Impl.Get()
	return nil
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a GreeterRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return GreeterRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type ShardPlugin struct {
	// Impl Injection
	Impl Shard
}

func (s *ShardPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ShardRPCServer{Impl: s.Impl}, nil
}

func (ShardPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ShardRPC{client: c}, nil
}
