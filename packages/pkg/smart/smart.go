/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package smart

import (
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"math/rand"
)

// SmartContract is storing smart contract data
type SmartContract struct {
	CLB             bool
	Rollback        bool
	FullAccess      bool
	SysUpdate       bool
	TxSmart         *types.SmartTransaction
	TxData          map[string]any
	TxFuel          int64           // The fuel of executing contract
	TxCost          int64           // Maximum cost of executing contract
	TxUsedCost      decimal.Decimal // Used cost of CPU resources
	TXBlockFuel     decimal.Decimal
	Loop            map[string]bool
	Hash            []byte
	Payload         []byte
	Timestamp       int64
	TxSignature     []byte
	TxSize          int64
	Size            common.StorageSize
	PublicKeys      [][]byte
	Rand            *rand.Rand
	Notifications   types.Notifications
	GenBlock        bool
	TimeLimit       int64
	taxes           bool
	Penalty         bool
	TokenEcosystems map[int64]any
	PrevSysPar      map[string]string
	ComPercents     map[int64]int64
}
