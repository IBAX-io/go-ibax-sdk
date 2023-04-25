package base

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/config"
	"github.com/IBAX-io/go-ibax-sdk/packages/modus"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/common/crypto"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/converter"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
	"github.com/IBAX-io/go-ibax-sdk/packages/response"
	"github.com/shopspring/decimal"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// MakeRequest builds a Request from all its parts, but returns an error if the
// params cannot be marshalled.
func MakeRequest(id uint64, jsonMethod string, params ...any) (request.Request, error) {
	p, err := request.MakeParams(params...)
	if err != nil {
		return request.Request{}, err
	}

	return request.Request{
		JSONRPC: request.JsonRPCVersion,
		ID:      request.ID{Num: id},
		Method:  jsonMethod,
		Params:  p,
	}, nil
}

// MustRequest builds a request from all its parts but panics if the params cannot be marshaled,
// so should only be used with well-known parameter data.
func MustRequest(id uint64, jsonMethod string, params ...any) (request.Request, error) {
	r, err := MakeRequest(id, jsonMethod, params...)
	if err != nil {
		return request.Request{}, err
	}

	return r, nil
}

// PrivateToPublicHex returns the hex public key for the specified hex private key.
func PrivateToPublicHex(hexkey string) (string, error) {
	key, err := hex.DecodeString(hexkey)
	if err != nil {
		return ``, fmt.Errorf("Decode hex error")
	}
	pubKey, err := crypto.PrivateToPublic(key)
	if err != nil {
		return ``, err
	}
	return crypto.PubToHex(pubKey), nil
}

type base struct {
	config *config.Config
	lock   sync.RWMutex
	id     uint64
}

func New(config config.Config) modus.Base {
	b := &base{config: &config}
	if err := b.Init(); err != nil {
		log.Fatalf("new base init failed:%s\n", err.Error())
	}
	return b
}

func (c *base) Init() (err error) {
	var (
		key, pub []byte
		pubStr   string
	)

	c.lock.Lock()
	defer c.lock.Unlock()
	crypto.InitAsymAlgo(c.config.Cryptoer)
	crypto.InitHashAlgo(c.config.Hasher)
	if c.config.PrivateKey == "" {
		return nil
	}
	key = []byte(c.config.PrivateKey)
	if len(key) > 64 {
		key = key[:64]
	}
	if len(key) != 64 {
		return fmt.Errorf("private key invalid")
	}
	c.config.PrivateKey = string(key)

	pubStr, err = PrivateToPublicHex(string(key))
	if err != nil {
		return fmt.Errorf("prase public key failed:%s", err.Error())
	}
	pub, err = hex.DecodeString(pubStr)
	if err != nil {
		return fmt.Errorf("invalid pubic key:%s", err.Error())
	}
	c.config.PublicKey = pub
	c.config.KeyId = crypto.Address(pub)
	c.config.Account = converter.AddressToString(c.config.KeyId)

	return
}

func (c *base) sendRequest(method string, form any, result any) error {
	return c.sendRawRequest(method, form, result)
}

func (c *base) ExpediteValidator(expedite string) error {
	//input Max unit expedite
	if expedite == "" {
		return errors.New("[expedite] is empty")
	}
	d, err := decimal.NewFromString(expedite)
	if err != nil {
		return fmt.Errorf("[expedite] params invalid:%s,err:%s", expedite, err.Error())
	}
	mix := decimal.New(1, -12)
	if d.LessThan(mix) || d.Mod(mix).GreaterThan(decimal.Zero) {
		return fmt.Errorf("[expedite] inconsistent with the smallest reference unit and its integer multiples:%s", expedite)
	}
	return nil
}

func (c *base) AmountValidator(amount string) error {
	if amount == "" {
		return errors.New("[amount] is empty")
	}
	d, err := decimal.NewFromString(amount)
	if err != nil {
		return fmt.Errorf("[amount] params invalid:%s,err:%s", amount, err.Error())
	}
	mix := decimal.New(1, 0)
	if d.LessThan(mix) || d.Mod(mix).GreaterThan(decimal.Zero) {
		return fmt.Errorf("[amount] inconsistent with the smallest reference unit and its integer multiples:%s", amount)
	}
	return nil
}

func (c *base) nextID() uint64 {
	return atomic.AddUint64(&c.id, 1)
}

func (c *base) newMethod(namespace request.Namespace, name string) (method string) {
	method = string(namespace) + request.NamespaceSeparator + name
	return
}

func (c *base) NewMessage(params request.RequestParams) (request.Request, error) {
	return MustRequest(c.nextID(), c.newMethod(params.Namespace, params.Name), params.Params...)
}

func (c *base) NewBatchMessage(requestPrams []request.BatchRequestParams) ([]request.BatchRequest, error) {
	var batchRequest []request.BatchRequest
	for _, params := range requestPrams {
		req, err := MustRequest(c.nextID(), c.newMethod(params.Namespace, params.Name), params.Params...)
		if err != nil {
			return nil, err
		}
		var bat request.BatchRequest
		bat.Result = params.Result
		bat.Req = &req
		batchRequest = append(batchRequest, bat)
	}
	return batchRequest, nil
}

var defaultClient = &http.Client{
	Timeout: 10 * time.Second,
}

func (c *base) sendRawRequest(method string, msg any, result any) error {
	var (
		isBatch bool
		body    []byte
		err     error
	)

	if msg != nil {
		switch msg.(type) {
		case request.Request:
		case []request.BatchRequest:
			isBatch = true
			reqs := msg.([]request.BatchRequest)
			batchRequest := make([]request.Request, len(reqs))
			for k, v := range reqs {
				batchRequest[k] = *v.Req
			}
			body, err = json.Marshal(batchRequest)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("the request structure does not support")
		}
	}
	if !isBatch {
		body, err = json.Marshal(msg)
		if err != nil {
			return err
		}
	}

	cnf := c.GetConfig()
	apiAddress := cnf.ApiAddress
	req, err := http.NewRequest(method, apiAddress, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.ContentLength = int64(len(body))

	req.Header.Set("Content-Type", "application/json")
	if len(cnf.Token) > 0 {
		req.Header.Set("Authorization", cnf.JwtPrefix+cnf.Token)
	}
	resp, err := defaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(`status code:%d`, resp.StatusCode)
	}
	if result == nil {
		return nil
	}

	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if isBatch {
			ret := &response.BatchResponse{}
			err = json.NewDecoder(resp.Body).Decode(ret)
			if err != nil {
				return fmt.Errorf("json response decode failed:%s", err.Error())
			}
			if len(*ret) == 0 {
				return errors.New("not result in JSON-RPC response")
			}

			bq := msg.([]request.BatchRequest)
			for _, v := range *ret {
				for k, b := range bq {
					if b.Req.ID == v.ID {
						if v.Error != nil {
							b.Err = err
							bq[k] = b
							break
						}
						if len(v.Result) == 0 {
							b.Err = errors.New("not result in JSON-RPC response")
							bq[k] = b
							break
						}

						err = json.Unmarshal(v.Result, &b.Result)
						if err != nil {
							b.Err = err
							bq[k] = b
							break
						}
						bq[k] = b
						break
					}
				}
			}

			return nil
		}
		ret := &response.Response{}
		err = json.NewDecoder(resp.Body).Decode(ret)
		if err != nil {
			return fmt.Errorf("json response decode failed:%s", err.Error())
		}

		if ret.Error != nil {
			return ret.Error
		}
		if len(ret.Result) == 0 {
			return errors.New("not result in JSON-RPC response")
		}

		err = json.Unmarshal(ret.Result, &result)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("response Content-Type not support:%s\n", contentType)
	}

	return nil
}

func (c *base) GET(form any, v any) error {
	return c.sendRequest(http.MethodGet, form, v)
}

func (c *base) POST(form any, v any) error {
	return c.sendRequest(http.MethodPost, form, v)
}

func (c *base) GetConfig() config.Config {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return *c.config
}

func (c *base) SetConfig(cnf config.Config) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.config = &cnf
}

func (c *base) Version() string {
	return config.Version
}

func (c *base) SendGet(url string, form *url.Values, result any) error {
	return response.NotSupportError
}
func (c *base) SendPost(url string, form *url.Values, result any) error {
	return response.NotSupportError
}
func (c *base) SendMultipart(url string, files map[string][]byte, result any) error {
	return response.NotSupportError
}
