package base

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/config"
	"github.com/IBAX-io/go-ibax-sdk/pkg/common/crypto"
	"github.com/IBAX-io/go-ibax-sdk/pkg/converter"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

const version = "1.0.0"

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

func (c *baseClient) Init() (err error) {
	var (
		key, pub []byte
		pubStr   string
	)

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

func (c *baseClient) sendRequest(method, url string, form *url.Values, v any) error {
	return c.sendRawRequest(method, url, form, v)
}

func (c *baseClient) sendRawRequest(method, url string, form *url.Values, v any) error {
	client := &http.Client{}
	var ioForm io.Reader
	if form != nil {
		ioForm = strings.NewReader(form.Encode())
	}
	apiAddress := c.config.ApiAddress + c.config.ApiPath
	req, err := http.NewRequest(method, apiAddress+url, ioForm)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if len(c.config.Token) > 0 {
		req.Header.Set("Authorization", c.config.JwtPrefix+c.config.Token)
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
		return fmt.Errorf(`%d %s`, resp.StatusCode, strings.TrimSpace(string(data)))
	}
	if v != nil {
		contentType := resp.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			err = json.Unmarshal(data, v)
			if err != nil {
				return fmt.Errorf("json unmarshal failed:%s", err.Error())
			}
		} else if strings.Contains(contentType, "application/octet-stream") {
			switch vt := v.(type) {
			case *string:
				*vt = string(data)
			case *[]byte:
				*vt = data
			case *interface{}:
				*vt = interface{}(data)
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

func (c *baseClient) SendMultipart(url string, files map[string][]byte, v any) error {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

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

	req, err := http.NewRequest("POST", c.config.ApiAddress+c.config.ApiPath+url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if len(c.config.Token) > 0 {
		req.Header.Set("Authorization", c.config.JwtPrefix+c.config.Token)
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
		return fmt.Errorf(`%d %s`, resp.StatusCode, strings.TrimSpace(string(data)))
	}

	return json.Unmarshal(data, &v)
}

func (c *baseClient) SendGet(url string, form *url.Values, v any) error {
	return c.sendRequest("GET", url, form, v)
}

func (c *baseClient) SendPost(url string, form *url.Values, v any) error {
	return c.sendRequest("POST", url, form, v)
}

func (c *baseClient) GetConfig() *config.IbaxConfig {
	c.config.RLock()
	defer c.config.RUnlock()
	return c.config
}

func (c *baseClient) SetConfig(cnf *config.IbaxConfig) {
	c.config.Lock()
	defer c.config.Unlock()
	c.config = cnf
}

func (c *baseClient) Version() string {
	return version
}
