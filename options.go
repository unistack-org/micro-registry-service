package service

import (
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/registry"
)

type clientKey struct{}

// Client sets the RPC client
func Client(c client.Client) registry.Option {
	return registry.SetOption(clientKey{}, c)
}
