package auth

import (
	"encoding/hex"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/pkg/common/crypto"
	"github.com/IBAX-io/go-ibax-sdk/pkg/converter"
	"github.com/IBAX-io/go-ibax-sdk/response"
	"net/url"
	"strconv"
	"time"
)

type getUIDResult struct {
	UID         string `json:"uid,omitempty"`
	Token       string `json:"token,omitempty"`
	Expire      string `json:"expire,omitempty"`
	EcosystemID string `json:"ecosystem_id,omitempty"`
	KeyID       string `json:"key_id,omitempty"`
	Address     string `json:"address,omitempty"`
	NetworkID   string `json:"network_id,omitempty"`
	Cryptoer    string `json:"cryptoer"`
	Hasher      string `json:"hasher"`
}

type loginResult struct {
	Token       string        `json:"token,omitempty"`
	EcosystemID string        `json:"ecosystem_id,omitempty"`
	KeyID       string        `json:"key_id,omitempty"`
	Account     string        `json:"account,omitempty"`
	NotifyKey   string        `json:"notify_key,omitempty"`
	IsNode      bool          `json:"isnode"`
	IsOwner     bool          `json:"isowner"`
	IsCLB       bool          `json:"clb"`
	Timestamp   string        `json:"timestamp,omitempty"`
	Roles       []rolesResult `json:"roles,omitempty"`
}

type rolesResult struct {
	RoleId   int64  `json:"role_id"`
	RoleName string `json:"role_name"`
}

func (c *authClient) GetUid() error {
	var ret getUIDResult

	err := c.baseClient.SendGet(`getuid`, nil, &ret)
	if err != nil {
		return err
	}
	if c.config.Token != "" && ret.Expire != "" {
		t1, err := time.ParseDuration(ret.Expire)
		if err != nil {
			return err
		}
		c.config.TokenExpireTime = time.Now().Add(t1).Unix()
		return nil
	}
	c.config.Token = ret.Token
	c.config.UID = ret.UID
	networkId, err := strconv.ParseInt(ret.NetworkID, 10, 64)
	if err != nil {
		return err
	}
	c.config.NetworkId = networkId
	if c.config.Cryptoer != ret.Cryptoer || c.config.Hasher != ret.Hasher {
		crypto.InitAsymAlgo(ret.Cryptoer)
		crypto.InitHashAlgo(ret.Hasher)
		err = c.baseClient.Init()
		if err != nil {
			return err
		}
	}

	return nil
}

// Login
// roleId : Get the role id through the keyinfo interface
func (c *authClient) Login(roleId string) (err error) {
	var (
		sign []byte
		ret  loginResult
	)

	sign, err = crypto.SignString(c.config.PrivateKey, "LOGIN"+strconv.FormatInt(c.config.NetworkId, 10)+c.config.UID)
	if err != nil {
		return
	}

	form := url.Values{"pubkey": {hex.EncodeToString(c.config.PublicKey)}, "signature": {hex.EncodeToString(sign)},
		`ecosystem`: {converter.Int64ToStr(c.config.Ecosystem)}}
	if roleId != "0" {
		form.Set("role_id", roleId)
	}
	if c.config.IsMobile {
		form[`mobile`] = []string{`true`}
	}
	err = c.baseClient.SendPost(`login`, &form, &ret)
	if err == nil {
		c.config.Token = ret.Token
	}
	return
}

func (c *authClient) AutoLogin() (err error) {
	if c.config.Token != "" {
		if time.Unix(c.config.TokenExpireTime, 0).Sub(time.Now()) < time.Minute*10 {
			c.config.Token = ""
		} else {
			return nil
		}
	}
	err = c.GetUid()
	if err == nil {
		//default 0
		err = c.Login("0")
		if err == nil {
			err = c.GetUid()
		}
	}
	return
}

func (c *authClient) GetUID() (*response.GetUIDResult, error) {
	var result response.GetUIDResult
	getUIdUrl := fmt.Sprintf("getuid")
	err := c.baseClient.SendGet(getUIdUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

// LOGIN
// roleId : Get the role id through the keyinfo interface
func (c *authClient) LOGIN(roleId string) (*response.LoginResult, error) {
	var (
		result response.LoginResult
	)
	loginUrl := fmt.Sprintf("login")

	sign, err := crypto.SignString(c.config.PrivateKey, "LOGIN"+strconv.FormatInt(c.config.NetworkId, 10)+c.config.UID)
	if err != nil {
		return &result, err
	}

	form := url.Values{"pubkey": {hex.EncodeToString(c.config.PublicKey)}, "signature": {hex.EncodeToString(sign)},
		`ecosystem`: {converter.Int64ToStr(c.config.Ecosystem)}}
	if roleId != "0" {
		form.Set("role_id", roleId)
	}
	if c.config.IsMobile {
		form[`mobile`] = []string{`true`}
	}

	err = c.baseClient.SendPost(loginUrl, &form, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}
