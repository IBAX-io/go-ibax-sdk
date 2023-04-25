package tx

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/packages/modus"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/converter"
	"github.com/IBAX-io/go-ibax-sdk/packages/response"
	"net/url"
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

type multiTxStatusResult struct {
	Results map[string]*txStatus `json:"results"`
}

type txStatusRequest struct {
	Hashes []string `json:"hashes"`
}

type tx struct {
	base modus.Base
}

func New(b modus.Base) modus.Transaction {
	return &tx{base: b}
}

// TxStatus
// hash: transaction hash
// interval: After the transaction is sent, the time interval for each query of the transaction status
func (c *tx) TxStatus(hash string, frequency int, interval time.Duration) (response.TxStatusResult, error) {
	return c.waitTx(hash, frequency, interval)
}

func (c *tx) waitTx(hash string, frequency int, interval time.Duration) (rets response.TxStatusResult, err error) {
	if interval.Milliseconds() < waitTxMinInterval.Milliseconds() {
		interval = waitTxMinInterval
	}
	time.Sleep(interval)
	data, err := json.Marshal(&txStatusRequest{
		Hashes: []string{hash},
	})
	if err != nil {
		return
	}
	rets.Hash = hash
	for i := 1; i < frequency; i++ {
		var multiRet multiTxStatusResult
		err = c.base.SendPost(`txstatus`, &url.Values{
			"data": {string(data)},
		}, &multiRet)
		if err != nil {
			return
		}

		ret := multiRet.Results[hash]
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
func (c *tx) TxsStatus(hashList []string, interval time.Duration) (map[string]response.TxStatusResult, error) {
	if interval.Milliseconds() < waitTxMinInterval.Milliseconds() {
		interval = waitTxMinInterval
	}
	time.Sleep(interval)
	data, err := json.Marshal(&txStatusRequest{
		Hashes: hashList,
	})
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
	var multiRet multiTxStatusResult
	err = c.base.SendPost("txstatus", &url.Values{
		"data": {string(data)},
	}, &multiRet)
	if err != nil {
		return nil, err
	}

	for hash, v := range multiRet.Results {
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

func (c *tx) SendTx(arrData map[string][]byte) (*map[string]string, error) {
	ret := &response.SendTxResult{}
	hashMap := map[string]string{}
	err := c.base.SendMultipart("sendTx", arrData, &ret)
	if err != nil {
		return &hashMap, err
	}
	hashMap = ret.Hashes
	return &hashMap, nil
}

func (c *tx) GetTxInfo(hash string, getContractInfo bool) (*response.TxInfoResult, error) {
	var result response.TxInfoResult
	reqUrl := fmt.Sprintf("txinfo/%s", hash)
	if getContractInfo {
		reqUrl += fmt.Sprintf("?contractinfo=1")
	}
	err := c.base.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}

	return &result, nil
}

func (c *tx) GetTxInfoMulti(hashList []string, getContractInfo bool) (*response.MultiTxInfoResult, error) {
	var result response.MultiTxInfoResult
	reqUrl := fmt.Sprintf("txinfomultiple")
	if len(hashList) == 0 {
		return &result, errors.New("params invalid")
	}
	data := strings.Join(hashList, ",")
	reqUrl += fmt.Sprintf("?data=%s", data)

	if getContractInfo {
		reqUrl += fmt.Sprintf("&contractinfo=1")
	}
	err := c.base.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}
