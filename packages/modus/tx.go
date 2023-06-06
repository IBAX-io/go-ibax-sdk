package modus

import (
	"github.com/IBAX-io/go-ibax-sdk/packages/response"
	"time"
)

// Transaction
// transaction
type Transaction interface {
	// SendTx
	// send transaction
	// return hashMap *map[string]string "transaction hash map"
	SendTx(arrData map[string][]byte) (hashMap *map[string]string, err error)
	// TxStatus
	// Query the transaction result of the specified hash
	TxStatus(hash string, frequency int, interval time.Duration) (response.TxStatusResult, error)
	// TxsStatus
	// Query the transaction result of the multiple hash
	TxsStatus(hashList []string, interval time.Duration) (map[string]response.TxStatusResult, error)

	GetTxInfo(hash string, getContractInfo bool) (*response.TxInfoResult, error)
	GetTxInfoMulti(hashList []string, getContractInfo bool) (*response.MultiTxInfoResult, error)
}

//
//type RpcTransaction interface {
//	Transaction
//}
//
//type ApiTransaction interface {
//	Transaction
//}
