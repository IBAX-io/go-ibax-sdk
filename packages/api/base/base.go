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
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

// PrivateToPublicHex returns the hex public key for the specified hex private key.
func PrivateToPublicHex(hexkey string) (string, error) {
	key, err := hex.DecodeString(hexkey)
	if err != nil {
		return ``, fmt.Errorf("decode hex error")
	}
	pubKey, err := crypto.PrivateToPublic(key)
	if err != nil {
		return ``, err
	}
	return crypto.PubToHex(pubKey), nil
}

type base struct {
	lock   sync.RWMutex
	config *config.Config
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
	if c.config.PrivateKey == "" {
		return nil
	}
	key = []byte(c.config.PrivateKey)
	if len(key) > 64 {
		key = key[:64]
	}
	c.config.PrivateKey = string(key)

	crypto.InitAsymAlgo(c.config.Cryptoer)
	crypto.InitHashAlgo(c.config.Hasher)

	pubStr, err = PrivateToPublicHex(string(key))
	if err != nil {
		return
	}
	pub, err = hex.DecodeString(pubStr)
	c.config.PublicKey = pub
	c.config.KeyId = crypto.Address(pub)
	c.config.Account = converter.AddressToString(c.config.KeyId)

	return
}

func (c *base) sendRequest(method, url string, form *url.Values, result any) error {
	return c.sendRawRequest(method, url, form, result)
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

func (c *base) sendRawRequest(method, url string, form *url.Values, result any) error {
	client := &http.Client{}
	var ioForm io.Reader
	if form != nil {
		ioForm = strings.NewReader(form.Encode())
	}
	cnf := c.GetConfig()
	apiAddress := cnf.ApiAddress + cnf.ApiPath
	req, err := http.NewRequest(method, apiAddress+url, ioForm)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if len(cnf.Token) > 0 {
		req.Header.Set("Authorization", cnf.JwtPrefix+cnf.Token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			c.config.Token = ""
		}
		return fmt.Errorf(`%d %s`, resp.StatusCode, strings.TrimSpace(string(data)))
	}
	if result != nil {
		contentType := resp.Header.Get("Content-Type")
		var isFileType bool
		switch result.(type) {
		case *request.FileType:
			isFileType = true
		}
		if strings.Contains(contentType, "application/json") && !isFileType {
			err = json.Unmarshal(data, result)
			if err != nil {
				return fmt.Errorf("json unmarshal failed:%s", err.Error())
			}
		} else {
			switch vt := result.(type) {
			case *string:
				*vt = string(data)
			case *[]byte:
				*vt = data
			case *interface{}:
				*vt = interface{}(data)
			case *request.FileType:
				fileName := result.(*request.FileType).Name
				if fileName != "" {
					err = os.WriteFile(fileName, data, 0644)
					if err != nil {
						return err
					}
					*vt = request.FileType{fileName, contentType, ""}
				} else {
					*vt = request.FileType{"", contentType, string(data)}
				}
			default:
				return fmt.Errorf("not supported TYPE:%T", vt)
			}

			//value = string(data)
			//value1 := reflect.ValueOf(value).Elem()
			//value1.SetString(string(data))
		}
	}

	return nil
}

func (c *base) SendMultipart(url string, files map[string][]byte, result any) error {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	cnf := c.GetConfig()

	for key, data := range files {
		part, err := writer.CreateFormFile(key, key)
		if err != nil {
			return err
		}
		if _, err := part.Write(data); err != nil {
			return err
		}
	}

	if err := writer.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", cnf.ApiAddress+cnf.ApiPath+url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if len(cnf.Token) > 0 {
		req.Header.Set("Authorization", cnf.JwtPrefix+cnf.Token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			c.config.Token = ""
		}
		return fmt.Errorf(`%d %s`, resp.StatusCode, strings.TrimSpace(string(data)))
	}

	return json.Unmarshal(data, &result)
}

func (c *base) SendGet(url string, form *url.Values, result any) error {
	return c.sendRequest("GET", url, form, result)
}

func (c *base) SendPost(url string, form *url.Values, result any) error {
	return c.sendRequest("POST", url, form, result)
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

func (c *base) GET(form any, result any) error {
	return response.NotSupportError
}

func (c *base) POST(form any, result any) error {
	return response.NotSupportError
}

func (c *base) NewMessage(params request.RequestParams) (request.Request, error) {
	return request.Request{}, response.NotSupportError
}

func (c *base) NewBatchMessage(requestPrams []request.BatchRequestParams) ([]request.BatchRequest, error) {
	return []request.BatchRequest{}, response.NotSupportError
}
