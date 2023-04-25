/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/
package types

import (
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/common/crypto"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/consts"
	"regexp"
	"strings"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
)

// Transaction types.
const (
	FirstBlockTxType = iota + 1
	StopNetworkTxType
	SmartContractTxType
	DelayTxType
	UtxoTxType
	TransferSelfTxType
)

// Header is contain header data
type Header struct {
	ID          int
	EcosystemID int64
	KeyID       int64
	Time        int64
	NetworkID   int64
	PublicKey   []byte
}

type TransferSelf struct {
	Value  string
	Source string
	Target string
}

// UTXO Transfer
type UTXO struct {
	ToID    int64
	Value   string
	Comment string
}

// SmartTransaction is storing smart contract data
type SmartTransaction struct {
	*Header
	MaxSum       string
	PayOver      string
	Lang         string
	Expedite     string
	SignedBy     int64
	TransferSelf *TransferSelf
	UTXO         *UTXO
	Params       map[string]any
}

func (s *SmartTransaction) TxType() byte {
	if s.TransferSelf != nil {
		return TransferSelfTxType
	}
	if s.UTXO != nil {
		return UtxoTxType
	}
	return SmartContractTxType
}

func (s *SmartTransaction) WithPrivate(privateKey []byte, internal bool) error {
	var (
		publicKey []byte
		err       error
	)
	if publicKey, err = crypto.PrivateToPublic(privateKey); err != nil {
		log.WithFields(log.Fields{"type": consts.CryptoError, "error": err}).Error("converting node private key to public")
		return err
	}
	s.PublicKey = publicKey
	if internal {
		s.SignedBy = crypto.Address(publicKey)
	}
	return nil
}

func (s *SmartTransaction) Unmarshal(buffer []byte) error {
	return msgpack.Unmarshal(buffer, s)
}

func (s *SmartTransaction) Marshal() ([]byte, error) {
	return msgpack.Marshal(s)
}

func (t SmartTransaction) Hash() ([]byte, error) {
	b, err := t.Marshal()
	if err != nil {
		return nil, err
	}
	return crypto.DoubleHash(b), nil
}

func (txSmart *SmartTransaction) Validate() error {
	if len(txSmart.Expedite) > 0 {
		expedite, _ := decimal.NewFromString(txSmart.Expedite)
		if expedite.LessThan(decimal.Zero) {
			return fmt.Errorf("expedite fee %s must be greater than 0", expedite)
		}
	}
	if len(strings.TrimSpace(txSmart.Lang)) > 2 {
		return fmt.Errorf(`localization size is greater than 2`)
	}
	if txSmart.UTXO != nil && len(txSmart.UTXO.Value) > 0 {
		if ok, _ := regexp.MatchString("^\\d+$", txSmart.UTXO.Value); !ok {
			return fmt.Errorf("error UTXO %s must integer", txSmart.UTXO.Value)
		}
	}
	if txSmart.TransferSelf != nil && len(txSmart.TransferSelf.Value) > 0 {
		if ok, _ := regexp.MatchString("^\\d+$", txSmart.TransferSelf.Value); !ok {
			return fmt.Errorf("error TransferSelf %s must integer", txSmart.TransferSelf.Value)
		}
	}

	return nil
}
