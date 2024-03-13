/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/
package transaction

import (
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/consts"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/converter"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/smart"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/types"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/utils"
	"github.com/IBAX-io/go-ibax/packages/common/crypto"
	"github.com/vmihailenco/msgpack/v5"
	"time"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type SmartTransactionParser struct {
	*smart.SmartContract
}

func (s *SmartTransactionParser) txType() byte      { return s.TxSmart.TxType() }
func (s *SmartTransactionParser) txHash() []byte    { return s.Hash }
func (s *SmartTransactionParser) txPayload() []byte { return s.Payload }
func (s *SmartTransactionParser) txTime() int64     { return s.Timestamp }
func (s *SmartTransactionParser) txKeyID() int64    { return s.TxSmart.KeyID }
func (s *SmartTransactionParser) txExpedite() decimal.Decimal {
	dec, _ := decimal.NewFromString(s.TxSmart.Expedite)
	return dec
}
func (s *SmartTransactionParser) setTimestamp() {
	s.Timestamp = time.Now().UnixMilli()
}

func (s *SmartTransactionParser) Validate() error {
	if err := s.TxSmart.Validate(); err != nil {
		return err
	}
	_, err := utils.CheckSign([][]byte{crypto.CutPub(s.TxSmart.PublicKey)}, s.Hash, s.TxSignature, false)
	if err != nil {
		return err
	}
	return nil
}

func (s *SmartTransactionParser) Marshal() ([]byte, error) {
	s.setTimestamp()
	if err := s.Validate(); err != nil {
		return nil, err
	}
	buf, err := msgpack.Marshal(s)
	if err != nil {
		return nil, err
	}
	buf = append([]byte{s.txType()}, buf...)
	return buf, nil
}

func (s *SmartTransactionParser) setSig(privateKey []byte) error {
	signature, err := crypto.Sign(privateKey, s.Hash)
	if err != nil {
		log.WithFields(log.Fields{"type": consts.CryptoError, "error": err}).Error("signing by node private key")
		return err
	}
	s.TxSignature = converter.EncodeLengthPlusData(signature)
	return nil
}

func (s *SmartTransactionParser) BinMarshalWithPrivate(smartTx *types.SmartTransaction, privateKey []byte, internal bool) ([]byte, error) {
	var (
		buf []byte
		err error
	)
	if err = smartTx.WithPrivate(privateKey, internal); err != nil {
		return nil, err
	}
	s.TxSmart = smartTx
	buf, err = s.TxSmart.Marshal()
	if err != nil {
		return nil, err
	}
	s.Payload = buf
	s.Hash = crypto.DoubleHash(s.Payload)
	err = s.setSig(privateKey)
	if err != nil {
		return nil, err
	}
	return s.Marshal()
}
