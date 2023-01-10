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
}

func amountValidator(value string) error {
	if value == "" {
		return nil
	}
	d, err := decimal.NewFromString(value)
	if err != nil {
		return fmt.Errorf("params invalid:%s,err:%s", value, err.Error())
	}
	if d.GreaterThanOrEqual(decimal.Zero) {
		return fmt.Errorf("params invalid:%s", value)
	}
	return nil
}

func (chain *utxoClient) NewUtxoSmartTransaction(txType utxoType, form getter, expedite string) (*types.SmartTransaction, error) {
	if expedite != "" {
		err := amountValidator(expedite)
		if err != nil {
			return &types.SmartTransaction{}, err
		}
	}
	smartTx := types.SmartTransaction{
		Header: &types.Header{
			Time:        time.Now().Unix(),
			EcosystemID: chain.config.Ecosystem,
			KeyID:       converter.StringToAddress(chain.config.Account),
			NetworkID:   chain.config.NetworkId,
		},
		Expedite: expedite,
	}
	amount := form.Get("amount")
	if len(amount) == 0 {
		return &smartTx, errors.New("amount params invalid")
	}
	err := amountValidator(amount)
	if err != nil {
		return &types.SmartTransaction{}, err
	}
	switch txType {
	case TypeTransfer:
		recipient := form.Get("recipient")
		comment := form.Get("comment")
		if len(recipient) == 0 {
			return &smartTx, errors.New("recipient params invalid")
		}
		toId := converter.StringToAddress(recipient)
		if toId == 0 && recipient != BlackHoleAddr {
			return &smartTx, fmt.Errorf("recipient %s is not valid", recipient)
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
	if privateKey, err = hex.DecodeString(chain.config.PrivateKey); err != nil {
		return
	}

	data, hash, err = transaction.NewTransactionInProc(smartTransaction, privateKey)

	return
}

func (chain *utxoClient) AutoCallUtxo(txType utxoType, form getter, expedite string) (*response.TxStatusResult, error) {
	var (
		rets = response.TxStatusResult{}
	)
	smartTx, err := chain.NewUtxoSmartTransaction(txType, form, expedite)
	if err != nil {
		return &rets, err
	}

	var privateKey []byte
	if privateKey, err = hex.DecodeString(chain.config.PrivateKey); err != nil {
		return &rets, err
	}

	arrData := make(map[string][]byte)
	data, txhash, err := transaction.NewTransactionInProc(*smartTx, privateKey)
	if err != nil {
		return &rets, err
	}
	arrData[fmt.Sprintf("%x", txhash)] = data
	//fmt.Println(fmt.Sprintf("%x", txhash))

	ret := &response.SendTxResult{}
	err = chain.SendMultipart("sendTx", arrData, &ret)
	if err != nil {
		return &rets, err
	}

	if len(form.Get("nowait")) > 0 {
		return &rets, nil
	}

	rets, err = chain.TxStatus(hex.EncodeToString(txhash), 10, time.Second*1)
	if err != nil {
		return &rets, err
	}
	return &rets, nil
}
