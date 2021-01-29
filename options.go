package service

import (
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/register"
)

type clientKey struct{}

// Client sets the RPC client
func Client(c client.Client) register.Option {
	return register.SetOption(clientKey{}, c)
}
