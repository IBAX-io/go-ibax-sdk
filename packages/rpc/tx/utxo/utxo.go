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

func (ux *utxo) NewUtxoSmartTransaction(txType request.UtxoType, form modus.Getter, expedite string) (*types.SmartTransaction, error) {
	if expedite != "" {
		//Uniform use min uint
		d, err := decimal.NewFromString(expedite)
		if err != nil {
			return &types.SmartTransaction{}, fmt.Errorf("expedite invalid:%s,err:%s", expedite, err.Error())
		}
		expedite = decimal.New(d.IntPart(), -12).String()
		err = ux.ExpediteValidator(expedite)
		if err != nil {
			return &types.SmartTransaction{}, err
		}
	}
	cnf := ux.GetConfig()
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
	err := ux.AmountValidator(amount)
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

func (ux *utxo) NewUtxoTransaction(smartTransaction types.SmartTransaction) (data, hash []byte, err error) {
	var privateKey []byte
	cnf := ux.GetConfig()
	if privateKey, err = hex.DecodeString(cnf.PrivateKey); err != nil {
		return
	}

	data, hash, err = transaction.NewTransactionInProc(smartTransaction, privateKey)

	return
}

func (ux *utxo) AutoCallUtxo(txType request.UtxoType, form modus.Getter, expedite string) (*response.TxStatusResult, error) {
	var (
		rets = response.TxStatusResult{}
	)
	smartTx, err := ux.NewUtxoSmartTransaction(txType, form, expedite)
	if err != nil {
		return &rets, err
	}

	cnf := ux.GetConfig()
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

	ret := &response.SendTxResult{}
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "sendTx",
		Params:    []any{arrData},
	}
	req, err := ux.NewMessage(message)
	if err != nil {
		return &rets, err
	}
	err = ux.POST(req, &ret)
	if err != nil {
		return &rets, err
	}

	if form != nil && len(form.Get("nowait")) > 0 {
		return &rets, nil
	}

	rets, err = ux.TxStatus(hex.EncodeToString(txhash), 10, time.Second*1)
	if err != nil {
		return &rets, err
	}
	return &rets, nil
}
