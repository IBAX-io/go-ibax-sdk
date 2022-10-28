package tx

import (
	"github.com/IBAX-io/go-ibax-sdk/base"
	response2 "github.com/IBAX-io/go-ibax-sdk/response"
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
	TxStatus(hash string, frequency int, interval time.Duration) (response2.TxStatusResult, error)
	// TxsStatus
	// Query the transaction result of the multiple hash
	TxsStatus(hashList []string, interval time.Duration) (map[string]response2.TxStatusResult, error)

	GetTxInfo(hash string, getContractInfo bool) (*response2.TxInfoResult, error)
	GetTxInfoMulti(hashList []string, getContractInfo bool) (*response2.MultiTxInfoResult, error)
}

type txClient struct {
	baseClient base.Base
}

func NewClient(b base.Base) Transaction {
	return &txClient{baseClient: b}
}
