package utxo

import (
	"github.com/IBAX-io/go-ibax-sdk/base"
	"github.com/IBAX-io/go-ibax-sdk/config"
	"github.com/IBAX-io/go-ibax-sdk/pkg/types"
	"github.com/IBAX-io/go-ibax-sdk/response"
	"github.com/IBAX-io/go-ibax-sdk/tx"
)

// Utxo
// Functions for Work with UTXO
type Utxo interface {
	NewUtxoSmartTransaction(txType utxoType, form getter, expedite string) (*types.SmartTransaction, error)
	NewUtxoTransaction(smartTransaction types.SmartTransaction) (data, hash []byte, err error)
	AutoCallUtxo(txType utxoType, form getter, expedite string) (*response.TxStatusResult, error)
}

type utxoClient struct {
	Config *config.IbaxConfig `yaml:"chain_sdk"`
	base.Base
	tx.Transaction
}

func NewClient(config *config.IbaxConfig, b base.Base, tx tx.Transaction) Utxo {
	return &utxoClient{Config: config, Base: b, Transaction: tx}
}
