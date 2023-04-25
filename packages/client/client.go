package client

import (
	"github.com/IBAX-io/go-ibax-sdk/config"
	"github.com/IBAX-io/go-ibax-sdk/packages/api"
	"github.com/IBAX-io/go-ibax-sdk/packages/modus"
	"github.com/IBAX-io/go-ibax-sdk/packages/rpc"
)

func NewClient(config config.Config) modus.Client {
	if config.EnableRpc {
		return rpc.NewClient(config)
	} else {
		return api.NewClient(config)
	}
}
