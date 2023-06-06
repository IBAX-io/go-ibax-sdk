package utxo

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/packages/modus"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/converter"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/transaction"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/types"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
	"github.com/IBAX-io/go-ibax-sdk/packages/response"
	"github.com/shopspring/decimal"
	"time"
)

type utxo struct {
	modus.Base
	modus.Transaction
}

func New(b modus.Base, tx modus.Transaction) modus.Utxo {
	return &utxo{Base: b, Transaction: tx}
}

func (u *utxo) NewUtxoSmartTransaction(txType request.UtxoType, form modus.Getter, expedite string) (*types.SmartTransaction, error) {
	if expedite != "" {
		//Uniform use min uint
		d, err := decimal.NewFromString(expedite)
		if err != nil {
			return &types.SmartTransaction{}, fmt.Errorf("expedite invalid:%s,err:%s", expedite, err.Error())
		}
		expedite = decimal.New(d.IntPart(), -12).String()
		err = u.ExpediteValidator(expedite)
		if err != nil {
			return &types.SmartTransaction{}, err
		}
	}
	cnf := u.GetConfig()
	smartTx := types.SmartTransaction{
		Header: &types.Header{
			Time:        time.Now().Unix(),
			EcosystemID: cnf.Ecosystem,
			KeyID:       converter.StringToAddress(cnf.Account),
			NetworkID:   cnf.NetworkId,
		},
		Expedite: expedite,
	}
	amount := form.Get("amount")
	if len(amount) == 0 {
		return &smartTx, errors.New("amount params invalid")
	}
	err := u.AmountValidator(amount)
	if err != nil {
		return &types.SmartTransaction{}, err
	}
	switch txType {
	case request.TypeTransfer:
		recipient := form.Get("recipient")
		comment := form.Get("comment")
		if len(recipient) == 0 {
			return &smartTx, errors.New("recipient params invalid")
		}
		toId := converter.StringToAddress(recipient)
		if toId == 0 && recipient != request.BlackHoleAddr {
			return &smartTx, fmt.Errorf("recipient %s is not valid", recipient)
		}
		smartTx.UTXO = &types.UTXO{
			Value:   amount,
			ToID:    toId,
			Comment: comment,
		}
	case request.TypeContractToUTXO:
		smartTx.TransferSelf = &types.TransferSelf{
			Value:  amount,
			Source: request.DestinationAccount,
			Target: request.DestinationUtxo,
		}
	case request.TypeUTXOToContract:
		smartTx.TransferSelf = &types.TransferSelf{
			Value:  amount,
			Source: request.DestinationUtxo,
			Target: request.DestinationAccount,
		}
	}

	return &smartTx, nil
}

func (u *utxo) NewUtxoTransaction(smartTransaction types.SmartTransaction) (data, hash []byte, err error) {
	var privateKey []byte
	cnf := u.GetConfig()
	if privateKey, err = hex.DecodeString(cnf.PrivateKey); err != nil {
		return
	}

	data, hash, err = transaction.NewTransactionInProc(smartTransaction, privateKey)

	return
}

func (u *utxo) AutoCallUtxo(txType request.UtxoType, form modus.Getter, expedite string) (*response.TxStatusResult, error) {
	var (
		rets = response.TxStatusResult{}
	)
	smartTx, err := u.NewUtxoSmartTransaction(txType, form, expedite)
	if err != nil {
		return &rets, err
	}

	cnf := u.GetConfig()
	var privateKey []byte
	if privateKey, err = hex.DecodeString(cnf.PrivateKey); err != nil {
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
	err = u.SendMultipart("sendTx", arrData, &ret)
	if err != nil {
		return &rets, err
	}

	if form != nil && len(form.Get("nowait")) > 0 {
		return &rets, nil
	}

	rets, err = u.TxStatus(hex.EncodeToString(txhash), 10, time.Second*1)
	if err != nil {
		return &rets, err
	}
	return &rets, nil
}
