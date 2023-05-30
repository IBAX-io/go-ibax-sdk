package contract

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/packages/modus"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/common/crypto"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/transaction"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/types"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
	"github.com/IBAX-io/go-ibax-sdk/packages/response"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

type contractField struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Optional bool   `json:"optional"`
}

type getContractInfo struct {
	ID       uint32          `json:"id"`
	StateID  uint32          `json:"state"`
	TableID  string          `json:"tableid"`
	WalletID string          `json:"walletid"`
	TokenID  string          `json:"tokenid"`
	Address  string          `json:"address"`
	Fields   []contractField `json:"fields"`
	Name     string          `json:"name"`
}

type contract struct {
	modus.Base
	modus.Transaction
}

func New(b modus.Base, tx modus.Transaction) modus.Contract {
	return &contract{Base: b, Transaction: tx}
}

// GetContracts
// Get Login ecosystem contracts
func (c *contract) GetContracts(limit, offset int64) (*response.ListResult, error) {
	var result response.ListResult
	getContractsUrl := fmt.Sprintf("contracts?limit=%d&offset=%d", limit, offset)
	err := c.SendGet(getContractsUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

// GetContract
// Get Login ecosystem contract by contract name
func (c *contract) GetContract(contractName string) (*response.GetContractResult, error) {
	var result response.GetContractResult
	getContractUrl := fmt.Sprintf("contract/%s", contractName)
	err := c.SendGet(getContractUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (c *contract) PrepareContractTx(contractName string, form modus.Getter) (params map[string]any, contractId int, err error) {
	var contract getContractInfo
	if err = c.SendGet("contract/"+contractName, nil, &contract); err != nil {
		return
	}

	params = make(map[string]any)
	for _, field := range contract.Fields {
		name := field.Name
		value := form.Get(name)

		if len(value) == 0 {
			continue
		}

		switch field.Type {
		case "bool":
			params[name], err = strconv.ParseBool(value)
		case "int", "address":
			params[name], err = strconv.ParseInt(value, 10, 64)
		case "float":
			params[name], err = strconv.ParseFloat(value, 64)
		case "map":
			var v map[string]any
			err = json.Unmarshal([]byte(value), &v)
			params[name] = v
		case "string", "money":
			params[name] = value
		case "file", "bytes", "array":
			if cp, ok := form.(*request.MapParams); !ok {
				err = fmt.Errorf("Form is not *contractParams type")
			} else {
				params[name] = cp.GetRaw(name)
			}
		}

		if err != nil {
			err = fmt.Errorf("parse param '%s':  %s,value:%s", name, value, err.Error())
			return
		}
	}
	return params, int(contract.ID), nil
}

func (c *contract) NewContractTransaction(contractId int, params map[string]any, expedite string) (data, hash []byte, err error) {
	if expedite != "" {
		//Uniform use min uint
		d, err := decimal.NewFromString(expedite)
		if err != nil {
			return nil, nil, fmt.Errorf("expedite invalid:%s,err:%s", expedite, err.Error())
		}
		expedite = decimal.New(d.IntPart(), -12).String()
		err = c.ExpediteValidator(expedite)
		if err != nil {
			return nil, nil, err
		}
	}
	cnf := c.Base.GetConfig()
	var privateKey, publicKey []byte
	if privateKey, err = hex.DecodeString(cnf.PrivateKey); err != nil {
		return
	}
	if publicKey, err = crypto.PrivateToPublic(privateKey); err != nil {
		return
	}
	data, hash, err = transaction.NewTransactionInProc(types.SmartTransaction{
		Header: &types.Header{
			ID:          contractId,
			Time:        time.Now().Unix(),
			EcosystemID: cnf.Ecosystem,
			KeyID:       crypto.Address(publicKey),
			NetworkID:   cnf.NetworkId,
		},
		Params:   params,
		Expedite: expedite,
	}, privateKey)

	return
}

func (c *contract) AutoCallContract(contractName string, form modus.Getter, expedite string) (*response.TxStatusResult, error) {
	var rets = response.TxStatusResult{}
	if expedite != "" {
		//Uniform use min uint
		d, err := decimal.NewFromString(expedite)
		if err != nil {
			return &rets, fmt.Errorf("expedite invalid:%s,err:%s", expedite, err.Error())
		}
		expedite = decimal.New(d.IntPart(), -12).String()
		err = c.ExpediteValidator(expedite)
		if err != nil {
			return &rets, err
		}
	}
	params, contractId, err := c.PrepareContractTx(contractName, form)
	if err != nil {
		return &rets, err
	}

	cnf := c.Base.GetConfig()
	var privateKey []byte
	if privateKey, err = hex.DecodeString(cnf.PrivateKey); err != nil {
		return &rets, err
	}

	arrData := make(map[string][]byte)
	data, txhash, err := transaction.NewTransactionInProc(types.SmartTransaction{
		Header: &types.Header{
			ID:          contractId,
			Time:        time.Now().Unix(),
			EcosystemID: cnf.Ecosystem,
			KeyID:       crypto.Address(cnf.PublicKey),
			NetworkID:   cnf.NetworkId,
		},
		Params:   params,
		Expedite: expedite,
	}, privateKey)
	if err != nil {
		return &rets, err
	}
	arrData[fmt.Sprintf("%x", txhash)] = data

	ret := &response.SendTxResult{}
	err = c.SendMultipart("sendTx", arrData, &ret)
	if err != nil {
		return &rets, err
	}

	if form != nil && len(form.Get("nowait")) > 0 {
		return &rets, nil
	}

	rets, err = c.TxStatus(hex.EncodeToString(txhash), 10, time.Second*1)
	if err != nil {
		return &rets, err
	}
	return &rets, nil
}
