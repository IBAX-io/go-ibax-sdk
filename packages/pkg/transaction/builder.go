/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/
package transaction

import (
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/consts"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/smart"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/types"
	log "github.com/sirupsen/logrus"
)

func newTransaction(smartTx types.SmartTransaction, privateKey []byte, internal bool) (data, hash []byte, err error) {
	stp := &SmartTransactionParser{
		SmartContract: &smart.SmartContract{TxSmart: new(types.SmartTransaction)},
	}
	data, err = stp.BinMarshalWithPrivate(&smartTx, privateKey, internal)
	if err != nil {
		log.WithFields(log.Fields{"type": consts.MarshallingError, "error": err}).Error("marshalling smart contract to msgpack")
		return
	}
	hash = stp.Hash
	return

}

func NewInternalTransaction(smartTx types.SmartTransaction, privateKey []byte) (data, hash []byte, err error) {
	return newTransaction(smartTx, privateKey, true)
}

func NewTransactionInProc(smartTx types.SmartTransaction, privateKey []byte) (data, hash []byte, err error) {
	return newTransaction(smartTx, privateKey, false)
}
