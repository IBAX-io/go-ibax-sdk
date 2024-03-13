package auth

import (
	"encoding/hex"
	"github.com/IBAX-io/go-ibax-sdk/packages/modus"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/converter"
	"github.com/IBAX-io/go-ibax-sdk/packages/response"
	"github.com/IBAX-io/go-ibax/packages/common/crypto"
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

type auth struct {
	base modus.Base
}

func New(b modus.Base) modus.Authentication {
	return &auth{base: b}
}

func (c *auth) GetAuthStatus() (*response.AuthStatusResponse, error) {
	return nil, response.NotSupportError
}

func (c *auth) GetUid() error {
	var ret getUIDResult

	err := c.base.SendGet(`getuid`, nil, &ret)
	if err != nil {
		return err
	}
	cnf := c.base.GetConfig()
	if cnf.Token != "" && ret.Expire != "" {
		t1, err := time.ParseDuration(ret.Expire)
		if err != nil {
			return err
		}
		cnf.TokenExpireTime = time.Now().Add(t1).Unix()
		c.base.SetConfig(cnf)
		return nil
	}
	cnf.Token = ret.Token
	cnf.UID = ret.UID
	networkId, err := strconv.ParseInt(ret.NetworkID, 10, 64)
	if err != nil {
		return err
	}
	cnf.NetworkId = networkId
	c.base.SetConfig(cnf)
	//Modify the encryption algorithm to use the encryption algorithm on the chain
	if cnf.Cryptoer != ret.Cryptoer || cnf.Hasher != ret.Hasher {
		cnf.Cryptoer = ret.Cryptoer
		cnf.Hasher = ret.Hasher
		c.base.SetConfig(cnf)
		err = c.base.Init()
		if err != nil {
			return err
		}
	}

	return nil
}

// Login
// roleId : Get the role id through the keyinfo interface
func (c *auth) Login(roleId int64) (err error) {
	var (
		sign []byte
		ret  loginResult
	)

	cnf := c.base.GetConfig()
	sign, err = crypto.SignString(cnf.PrivateKey, "LOGIN"+strconv.FormatInt(cnf.NetworkId, 10)+cnf.UID)
	if err != nil {
		return
	}

	form := url.Values{"pubkey": {hex.EncodeToString(cnf.PublicKey)}, "signature": {hex.EncodeToString(sign)},
		`ecosystem`: {converter.Int64ToStr(cnf.Ecosystem)}}
	if roleId != 0 {
		form.Set("role_id", strconv.FormatInt(roleId, 10))
	}
	err = c.base.SendPost(`login`, &form, &ret)
	if err == nil {
		cnf.Token = ret.Token
		c.base.SetConfig(cnf)
	}
	return
}

func (c *auth) AutoLogin() (err error) {
	cnf := c.base.GetConfig()
	if cnf.Token != "" {
		if time.Unix(cnf.TokenExpireTime, 0).Sub(time.Now()) < time.Minute*10 {
			cnf.Token = ""
			c.base.SetConfig(cnf)
		} else {
			return nil
		}
	}
	err = c.GetUid()
	if err == nil {
		//default 0
		err = c.Login(0)
		if err == nil {
			err = c.GetUid()
		}
	}
	return
}
