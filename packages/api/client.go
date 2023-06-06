package api

import (
	"github.com/IBAX-io/go-ibax-sdk/config"
	"github.com/IBAX-io/go-ibax-sdk/packages/api/auth"
	"github.com/IBAX-io/go-ibax-sdk/packages/api/base"
	"github.com/IBAX-io/go-ibax-sdk/packages/api/query"
	"github.com/IBAX-io/go-ibax-sdk/packages/api/tx"
	"github.com/IBAX-io/go-ibax-sdk/packages/api/tx/contract"
	"github.com/IBAX-io/go-ibax-sdk/packages/api/tx/utxo"
	"github.com/IBAX-io/go-ibax-sdk/packages/modus"
	"github.com/IBAX-io/go-ibax-sdk/packages/wallet"
)

type client struct {
	config *config.Config
	modus.Base
	modus.Authentication
	modus.Contract
	modus.Utxo
	modus.Transaction
	modus.Query
	modus.Wallet
}

func NewClient(config config.Config) modus.Client {
	b := base.New(config)
	a := auth.New(b)
	t := tx.New(b)
	c := contract.New(b, t)
	q := query.New(b)
	u := utxo.New(b, t)
	acc := wallet.New(b)
	return &client{Authentication: a, Base: b, Contract: c, Transaction: t, Query: q, Utxo: u, Wallet: acc}
}
