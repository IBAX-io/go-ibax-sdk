package contract

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/pkg/common/crypto"
	"github.com/IBAX-io/go-ibax-sdk/pkg/transaction"
	"github.com/IBAX-io/go-ibax-sdk/pkg/types"
	"github.com/IBAX-io/go-ibax-sdk/response"
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

type getter interface {
	Get(string) string
}

type ContractParams map[string]any

func (cp *ContractParams) Get(key string) string {
	if _, ok := (*cp)[key]; !ok {
		return ""
	}
	return fmt.Sprintf("%v", (*cp)[key])
}

func (cp *ContractParams) getRaw(key string) any {
	return (*cp)[key]
}

// GetContracts
// Get Login ecosystem contracts
func (c *contractClient) GetContracts(limit, offset int64) (*response.ListResult, error) {
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
func (c *contractClient) GetContract(contractName string) (*response.GetContractResult, error) {
	var result response.GetContractResult
	getContractUrl := fmt.Sprintf("contract/%s", contractName)
	err := c.SendGet(getContractUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (c *contractClient) PrepareContractTx(contractName string, form getter) (params map[string]any, contractId int, err error) {
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
		case "array":
			var v any
			err = json.Unmarshal([]byte(value), &v)
			params[name] = v
		case "map":
			var v map[string]any
			err = json.Unmarshal([]byte(value), &v)
			params[name] = v
		case "string", "money":
			params[name] = value
		case "file", "bytes":
			if cp, ok := form.(*ContractParams); !ok {
				err = fmt.Errorf("Form is not *contractParams type")
			} else {
				params[name] = cp.getRaw(name)
			}
		}

		if err != nil {
			err = fmt.Errorf("Parse param '%s': %s", name, err)
			return
		}
	}
	return params, int(contract.ID), nil
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

func (c *contractClient) NewContractTransaction(contractId int, params map[string]any, expedite string) (data, hash []byte, err error) {
	if expedite != "" {
		err = amountValidator(expedite)
		if err != nil {
			return nil, nil, err
		}
	}
	var privateKey, publicKey []byte
	if privateKey, err = hex.DecodeString(c.config.PrivateKey); err != nil {
		return
	}
	if publicKey, err = crypto.PrivateToPublic(privateKey); err != nil {
		return
	}
	data, hash, err = transaction.NewTransactionInProc(types.SmartTransaction{
		Header: &types.Header{
			ID:          contractId,
			Time:        time.Now().Unix(),
			EcosystemID: c.config.Ecosystem,
			KeyID:       crypto.Address(publicKey),
			NetworkID:   c.config.NetworkId,
		},
		Params:   params,
		Expedite: expedite,
	}, privateKey)

	return
}

func (c *contractClient) AutoCallContract(contractName string, form getter, expedite string) (*response.TxStatusResult, error) {
	var rets = response.TxStatusResult{}
	if expedite != "" {
		err := amountValidator(expedite)
		if err != nil {
			return &rets, err
		}
	}
	params, contractId, err := c.PrepareContractTx(contractName, form)
	if err != nil {
		return &rets, err
	}

	var privateKey []byte
	if privateKey, err = hex.DecodeString(c.config.PrivateKey); err != nil {
		return &rets, err
	}

	arrData := make(map[string][]byte)
	data, txhash, err := transaction.NewTransactionInProc(types.SmartTransaction{
		Header: &types.Header{
			ID:          contractId,
			Time:        time.Now().Unix(),
			EcosystemID: c.config.Ecosystem,
			KeyID:       crypto.Address(c.config.PublicKey),
			NetworkID:   c.config.NetworkId,
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

	if len(form.Get("nowait")) > 0 {
		return &rets, nil
	}

	rets, err = c.TxStatus(hex.EncodeToString(txhash), 10, time.Second*1)
	if err != nil {
		return &rets, err
	}
	return &rets, nil
}
