package tx

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/pkg/converter"
	response2 "github.com/IBAX-io/go-ibax-sdk/response"
	"net/url"
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

// TxStatus
// hash: transaction hash
// interval: After the transaction is sent, the time interval for each query of the transaction status
func (c *txClient) TxStatus(hash string, frequency int, interval time.Duration) (response2.TxStatusResult, error) {
	return c.waitTx(hash, frequency, interval)
}

func (c *txClient) waitTx(hash string, frequency int, interval time.Duration) (rets response2.TxStatusResult, err error) {
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

	for i := 1; i < frequency; i++ {
		var multiRet multiTxStatusResult
		err = c.baseClient.SendPost(`txstatus`, &url.Values{
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
					rets.Err = err
					return
				}
				rets.Err = errors.New(string(errText))
				return
			} else {
				rets.Err = fmt.Errorf(ret.Result)
				return
			}
		}
		if ret.Message != nil {
			errText, err = json.Marshal(ret.Message)
			if err != nil {
				rets.Err = err
				return
			}
			rets.Err = errors.New(string(errText))
			return
		}
		time.Sleep(interval)
	}

	return rets, nil
}

// TxsStatus
// hashList: multiple transaction hash
// interval: After the transaction is sent, the time interval for each query of the transaction status
func (c *txClient) TxsStatus(hashList []string, interval time.Duration) (map[string]response2.TxStatusResult, error) {
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
	rets := make(map[string]response2.TxStatusResult)
	setTxStatus := func(hash string, result response2.TxStatusResult) {
		rets[hash] = result
	}
	var againNumber int

	// again: If the transaction status is not queried, try again
again:
	var multiRet multiTxStatusResult
	err = c.baseClient.SendPost("txstatus", &url.Values{
		"data": {string(data)},
	}, &multiRet)
	if err != nil {
		return nil, err
	}

	for hash, v := range multiRet.Results {
		var result response2.TxStatusResult
		var errtext []byte
		if len(v.BlockID) > 0 {
			result.BlockId = converter.StrToInt64(v.BlockID)
			result.Penalty = v.Penalty
			if v.Penalty == 1 {
				errtext, err = json.Marshal(v.Message)
				if err != nil {
					result.Err = err
					setTxStatus(hash, result)
					continue
				}
				err = errors.New(string(errtext))
				result.Err = err
				setTxStatus(hash, result)
				continue
			} else {
				result.Err = fmt.Errorf(v.Result)
				setTxStatus(hash, result)
				continue
			}
		}
		if v.Message != nil {
			errtext, err = json.Marshal(v.Message)
			if err != nil {
				result.Err = fmt.Errorf(v.Result)
				setTxStatus(hash, result)
				continue
			}
			result.Err = errors.New(string(errtext))
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

func (c *txClient) SendTx(arrData map[string][]byte) (hashMap *map[string]string, err error) {
	ret := &response2.SendTxResult{}
	err = c.baseClient.SendMultipart("sendTx", arrData, &ret)
	if err != nil {
		return
	}
	hashMap = &ret.Hashes
	return
}

func (c *txClient) GetTxInfo(hash string, getContractInfo bool) (*response2.TxInfoResult, error) {
	var result response2.TxInfoResult
	reqUrl := fmt.Sprintf("txinfo/%s", hash)
	if getContractInfo {
		reqUrl += fmt.Sprintf("?contractinfo=1")
	}
	err := c.baseClient.SendGet(reqUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *txClient) GetTxInfoMulti(hashList []string, getContractInfo bool) (*response2.MultiTxInfoResult, error) {
	var result response2.MultiTxInfoResult
	reqUrl := fmt.Sprintf("txinfomultiple")
	if len(hashList) == 0 {
		return nil, errors.New("params invalid")
	}
	var request struct {
		Hashes []string `json:"hashes"`
	}
	request.Hashes = hashList
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	reqUrl += fmt.Sprintf("?data=%s", string(data))

	if getContractInfo {
		reqUrl += fmt.Sprintf("&contractinfo=1")
	}
	err = c.baseClient.SendGet(reqUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
