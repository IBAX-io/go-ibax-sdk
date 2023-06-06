package modus

import (
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/types"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
	"github.com/IBAX-io/go-ibax-sdk/packages/response"
)

// Utxo
// Functions for Work with UTXO
type Utxo interface {
	NewUtxoSmartTransaction(txType request.UtxoType, form Getter, expedite string) (*types.SmartTransaction, error)
	NewUtxoTransaction(smartTransaction types.SmartTransaction) (data, hash []byte, err error)
	AutoCallUtxo(txType request.UtxoType, form Getter, expedite string) (*response.TxStatusResult, error)
}
