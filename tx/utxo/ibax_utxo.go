package utxo

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/pkg/converter"
	"github.com/IBAX-io/go-ibax-sdk/pkg/transaction"
	"github.com/IBAX-io/go-ibax-sdk/pkg/types"
	"github.com/IBAX-io/go-ibax-sdk/response"
	"github.com/shopspring/decimal"
	"regexp"
	"time"
)

type utxoType int

const (
	TypeContractToUTXO utxoType = iota //Contract Account transfer UTXO Account
	TypeUTXOToContract                 //UTXO Account transfer Contract Account
	TypeTransfer
)

const (
	DestinationUtxo    = "UTXO"
	DestinationAccount = "Account"

	BlackHoleAddr = "0000-0000-0000-0000-0000"
)

type getter interface {
	Get(string) string
	Set(string, string)
}

func amountValidator(amount string) error {
	if ok, _ := regexp.MatchString("^\\d+$", amount); !ok {
		return errors.New("error new utxo smart transaction amount must be a positive integer")
	}
	if value, err := decimal.NewFromString(amount); err != nil || value.LessThanOrEqual(decimal.Zero) {
		return errors.New("error new utxo smart transaction amount must be greater than zero")
	}
	return nil
}

func (chain *utxoClient) NewUtxoSmartTransaction(txType utxoType, form getter, expedite string) (*types.SmartTransaction, error) {
	smartTx := types.SmartTransaction{
		Header: &types.Header{
			Time:        time.Now().Unix(),
			EcosystemID: chain.Config.Ecosystem,
			KeyID:       converter.StringToAddress(chain.Config.Account),
			NetworkID:   chain.Config.NetworkId,
		},
		Expedite: expedite,
	}
	amount := form.Get("amount")
	if len(amount) == 0 {
		return nil, errors.New("amount params invalid")
	}
	err := amountValidator(amount)
	if err != nil {
		return nil, err
	}
	switch txType {
	case TypeTransfer:
		recipient := form.Get("recipient")
		comment := form.Get("comment")
		if len(recipient) == 0 {
			return nil, errors.New("recipient params invalid")
		}
		toId := converter.StringToAddress(recipient)
		if toId == 0 && recipient != BlackHoleAddr {
			return nil, fmt.Errorf("recipient %s is not valid", recipient)
		}
		smartTx.UTXO = &types.UTXO{
			Value:   amount,
			ToID:    toId,
			Comment: comment,
		}
	case TypeContractToUTXO:
		smartTx.TransferSelf = &types.TransferSelf{
			Value:  amount,
			Source: DestinationAccount,
			Target: DestinationUtxo,
		}
	case TypeUTXOToContract:
		smartTx.TransferSelf = &types.TransferSelf{
			Value:  amount,
			Source: DestinationUtxo,
			Target: DestinationAccount,
		}
	}

	return &smartTx, nil
}

func (chain *utxoClient) NewUtxoTransaction(smartTransaction types.SmartTransaction) (data, hash []byte, err error) {
	var privateKey []byte
	if privateKey, err = hex.DecodeString(chain.Config.PrivateKey); err != nil {
		return
	}

	data, hash, err = transaction.NewTransactionInProc(smartTransaction, privateKey)

	return
}

func (chain *utxoClient) AutoCallUtxo(txType utxoType, form getter, expedite string) (*response.TxStatusResult, error) {
	smartTx, err := chain.NewUtxoSmartTransaction(txType, form, expedite)
	if err != nil {
		return nil, err
	}

	var privateKey []byte
	if privateKey, err = hex.DecodeString(chain.Config.PrivateKey); err != nil {
		return nil, err
	}

	arrData := make(map[string][]byte)
	data, txhash, err := transaction.NewTransactionInProc(*smartTx, privateKey)
	if err != nil {
		return nil, err
	}
	arrData[fmt.Sprintf("%x", txhash)] = data
	//fmt.Println(fmt.Sprintf("%x", txhash))

	ret := &response.SendTxResult{}
	err = chain.SendMultipart("sendTx", arrData, &ret)
	if err != nil {
		return nil, err
	}

	if len(form.Get("nowait")) > 0 {
		return nil, nil
	}

	rets, err := chain.TxStatus(hex.EncodeToString(txhash), 10, time.Second*1)
	if err != nil {
		return nil, err
	}
	return &rets, nil
}
