package contract

import (
	"github.com/IBAX-io/go-ibax-sdk/base"
	"github.com/IBAX-io/go-ibax-sdk/config"
	response2 "github.com/IBAX-io/go-ibax-sdk/response"
	"github.com/IBAX-io/go-ibax-sdk/tx"
)

// Contract
// Functions for Work with Contracts
type Contract interface {
	// GetContracts
	// @method get
	// get contracts
	GetContracts(limit, offset int64) (*response2.ListResult, error)
	// GetContract
	// @method get
	// get contract information by contract name
	GetContract(contractName string) (*response2.GetContractResult, error)
	// PrepareContractTx
	// @method get
	// contain Prepare Contract transaction function
	// @return params map[string]any "contract params"
	// @return contractId int "contract id"
	PrepareContractTx(contractName string, form getter) (params map[string]any, contractId int, err error)
	// NewContractTransaction
	// Build a contract transaction
	// @return data []byte "contract transaction data"
	// @return hash []byte "transaction data hash"
	NewContractTransaction(contractId int, params map[string]any, expedite string) (data, hash []byte, err error)
	// AutoCallContract
	// call Contract and return transaction result
	AutoCallContract(contractName string, form getter, expedite string) (*response2.TxStatusResult, error)
}

type contractClient struct {
	config *config.IbaxConfig
	base.Base
	tx.Transaction
}

func NewClient(config *config.IbaxConfig, b base.Base, tx tx.Transaction) Contract {
	return &contractClient{config: config, Base: b, Transaction: tx}
}
