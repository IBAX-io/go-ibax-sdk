package tx

import (
	"encoding/json"
	"errors"
	"github.com/IBAX-io/go-ibax-sdk/packages/modus"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/converter"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
	"github.com/IBAX-io/go-ibax-sdk/packages/response"
	"strconv"
	"strings"
	"time"
)

// After the transaction is sent, the min time interval for each query of the transaction status
const waitTxMinInterval = time.Millisecond * 500

type txstatusError struct {
	Type  string `json:"type,omitempty"`
	Error string `json:"error,omitempty"`
}

type txStatus struct {
	BlockID string         `json:"blockid"`
	Message *txstatusError `json:"errmsg,omitempty"`
	Result  string         `json:"result"`
	Penalty int64          `json:"penalty"`
}

type tx struct {
	baseClient modus.Base
}

func New(b modus.Base) modus.Transaction {
	return &tx{baseClient: b}
}

// TxStatus
// hash: transaction hash
// interval: After the transaction is sent, the time interval for each query of the transaction status
func (t *tx) TxStatus(hash string, frequency int, interval time.Duration) (response.TxStatusResult, error) {
	return t.waitTx(hash, frequency, interval)
}

func (t *tx) waitTx(hash string, frequency int, interval time.Duration) (rets response.TxStatusResult, err error) {
	if interval.Milliseconds() < waitTxMinInterval.Milliseconds() {
		interval = waitTxMinInterval
	}
	time.Sleep(interval)
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "txStatus",
		Params:    []any{hash},
	}
	req, err := t.baseClient.NewMessage(message)
	if err != nil {
		return rets, err
	}
	rets.Hash = hash
	for i := 1; i < frequency; i++ {
		var multiRet map[string]*txStatus
		err = t.baseClient.POST(req, &multiRet)
		if err != nil {
			return
		}
		if len(multiRet) == 0 {
			time.Sleep(interval)
			continue
		}

		ret := multiRet[hash]
		var errText []byte
		if len(ret.BlockID) > 0 {
			rets.BlockId = converter.StrToInt64(ret.BlockID)
			rets.Penalty = ret.Penalty
			if ret.Penalty == 1 {
				errText, err = json.Marshal(ret.Message)
				if err != nil {
					rets.Err = err.Error()
					return
				}
				rets.Err = string(errText)
				return
			} else {
				if ret.Result != "" {
					rets.Err = ret.Result
				}
				return
			}
		}
		if ret.Message != nil {
			errText, err = json.Marshal(ret.Message)
			if err != nil {
				rets.Err = err.Error()
				return
			}
			rets.Err = string(errText)
			return
		}
		time.Sleep(interval)
	}

	return rets, nil
}

// TxsStatus
// hashList: multiple transaction hash
// interval: After the transaction is sent, the time interval for each query of the transaction status
func (t *tx) TxsStatus(hashList []string, interval time.Duration) (map[string]response.TxStatusResult, error) {
	if interval.Milliseconds() < waitTxMinInterval.Milliseconds() {
		interval = waitTxMinInterval
	}
	time.Sleep(interval)
	hashes := strings.Join(hashList, ",")
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "txStatus",
		Params:    []any{hashes},
	}
	req, err := t.baseClient.NewMessage(message)
	if err != nil {
		return nil, err
	}

	rets := make(map[string]response.TxStatusResult)
	setTxStatus := func(hash string, result response.TxStatusResult) {
		result.Hash = hash
		rets[hash] = result
	}
	var againNumber int

	// again: If the transaction status is not queried, try again
again:
	//var multiRet multiTxStatusResult
	var multiRet map[string]*txStatus
	err = t.baseClient.POST(req, &multiRet)
	if err != nil {
		return nil, err
	}

	for hash, v := range multiRet {
		var result response.TxStatusResult
		var errtext []byte
		if len(v.BlockID) > 0 {
			result.BlockId = converter.StrToInt64(v.BlockID)
			result.Penalty = v.Penalty
			if v.Penalty == 1 {
				errtext, err = json.Marshal(v.Message)
				if err != nil {
					result.Err = err.Error()
					setTxStatus(hash, result)
					continue
				}
				err = errors.New(string(errtext))
				result.Err = err.Error()
				setTxStatus(hash, result)
				continue
			} else {
				result.Err = v.Result
				setTxStatus(hash, result)
				continue
			}
		}
		if v.Message != nil {
			errtext, err = json.Marshal(v.Message)
			if err != nil {
				result.Err = v.Result
				setTxStatus(hash, result)
				continue
			}
			result.Err = string(errtext)
			setTxStatus(hash, result)
			continue
		}

	}
	if len(rets) == 0 && againNumber < 10 {
		againNumber++
		time.Sleep(interval)
		goto again
	}

	return rets, nil
}

func (t *tx) SendTx(arrData map[string][]byte) (*map[string]string, error) {
	ret := &response.SendTxResult{}
	hashMap := map[string]string{}
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "sendTx",
		Params:    []any{arrData},
	}
	req, err := t.baseClient.NewMessage(message)
	if err != nil {
		return &hashMap, err
	}
	err = t.baseClient.POST(req, &ret)
	if err != nil {
		return &hashMap, err
	}

	hashMap = ret.Hashes
	return &hashMap, nil
}

type txInfoResult struct {
	BlockID int64            `json:"blockid"`
	Confirm int              `json:"confirm"`
	Data    *response.TxInfo `json:"data"`
}

func (t *tx) GetTxInfo(hash string, getContractInfo bool) (*response.TxInfoResult, error) {
	var rlt response.TxInfoResult
	var result txInfoResult
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "txInfo",
		Params:    []any{hash, getContractInfo},
	}
	req, err := t.baseClient.NewMessage(message)
	if err != nil {
		return &rlt, err
	}
	err = t.baseClient.POST(req, &result)
	if err != nil {
		return &rlt, err
	}
	rlt.Data = result.Data
	rlt.BlockID = strconv.FormatInt(result.BlockID, 10)
	rlt.Confirm = result.Confirm

	return &rlt, nil
}

func (t *tx) GetTxInfoMulti(hashList []string, getContractInfo bool) (*response.MultiTxInfoResult, error) {
	type multiTxInfoResult struct {
		Results map[string]*txInfoResult `json:"results"`
	}
	var result response.MultiTxInfoResult
	var rlt multiTxInfoResult
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "txInfoMultiple",
		Params:    []any{hashList, getContractInfo},
	}
	req, err := t.baseClient.NewMessage(message)
	if err != nil {
		return &result, err
	}

	err = t.baseClient.GET(req, &rlt)
	if err != nil {
		return &result, err
	}
	result.Results = map[string]*response.TxInfoResult{}
	for k, v := range rlt.Results {
		var rs response.TxInfoResult
		if v != nil {
			rs.Data = v.Data
			rs.Confirm = v.Confirm
			rs.BlockID = strconv.FormatInt(v.BlockID, 10)
		}
		result.Results[k] = &rs
	}

	return &result, nil
}
