package client

import (
	"github.com/IBAX-io/go-ibax-sdk/auth"
	"github.com/IBAX-io/go-ibax-sdk/base"
	"github.com/IBAX-io/go-ibax-sdk/config"
	"github.com/IBAX-io/go-ibax-sdk/query"
	"github.com/IBAX-io/go-ibax-sdk/tx"
	"github.com/IBAX-io/go-ibax-sdk/tx/contract"
	"github.com/IBAX-io/go-ibax-sdk/tx/utxo"
)

type client struct {
	config *config.IbaxConfig
	base.Base
	auth.Authentication
	contract.Contract
	utxo.Utxo
	tx.Transaction
	query.Query
}

type Client interface {
	base.Base
	auth.Authentication
	contract.Contract
	utxo.Utxo
	tx.Transaction
	query.Query
}

func NewClient(config *config.IbaxConfig) Client {
	b := base.NewClient(config)
	a := auth.NewClient(config, b)
	t := tx.NewClient(b)
	c := contract.NewClient(config, b, t)
	q := query.NewClient(b)
	u := utxo.NewClient(config, b, t)
	return &client{config: config, Authentication: a, Base: b, Contract: c, Transaction: t, Query: q, Utxo: u}
}
