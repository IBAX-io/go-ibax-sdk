package modus

import "github.com/IBAX-io/go-ibax-sdk/packages/response"

type Getter interface {
	Get(string) string
}

// Contract
// Functions for Work with Contracts
type Contract interface {
	// GetContracts
	// @method get
	// get contracts
	GetContracts(limit, offset int64) (*response.ListResult, error)
	// GetContract
	// @method get
	// get contract information by contract name
	GetContract(contractName string) (*response.GetContractResult, error)
	// PrepareContractTx
	// @method get
	// contain Prepare Contract transaction function
	// @return params map[string]any "contract params"
	// @return contractId int "contract id"
	PrepareContractTx(contractName string, form Getter) (params map[string]any, contractId uint32, err error)
	// NewContractTransaction
	// Build a contract transaction
	// @return data []byte "contract transaction data"
	// @return hash []byte "transaction data hash"
	// expedite: ibax fee
	NewContractTransaction(contractId uint32, params map[string]any, expedite string) (data, hash []byte, err error)
	// AutoCallContract
	// call Contract and return transaction result
	AutoCallContract(contractName string, form Getter, expedite string) (*response.TxStatusResult, error)
}
