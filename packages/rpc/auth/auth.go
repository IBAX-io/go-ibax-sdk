package auth

import (
	"encoding/hex"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/packages/modus"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
	"github.com/IBAX-io/go-ibax-sdk/packages/response"
	"github.com/IBAX-io/go-ibax/packages/common/crypto"
	"strconv"
	"time"
)

type auth struct {
	base modus.Base
}

func New(b modus.Base) modus.Authentication {
	return &auth{base: b}
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

func (a *auth) GetUidResponse() (*response.GetUIDResult, error) {
	var ret response.GetUIDResult

	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getUid",
	}
	req, err := a.base.NewMessage(message)
	if err != nil {
		return nil, err
	}
	err = a.base.GET(req, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (a *auth) GetUid() error {
	var ret response.GetUIDResult

	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getUid",
	}
	req, err := a.base.NewMessage(message)
	if err != nil {
		return err
	}
	err = a.base.GET(req, &ret)
	if err != nil {
		return err
	}
	cnf := a.base.GetConfig()

	if cnf.Token != "" && ret.Expire != "" {
		t1, err := time.ParseDuration(ret.Expire)
		if err != nil {
			return err
		}
		cnf.TokenExpireTime = time.Now().Add(t1).Unix()
		a.base.SetConfig(cnf)
		return nil
	}
	cnf.Token = ret.Token
	cnf.UID = ret.UID
	networkId, err := strconv.ParseInt(ret.NetworkID, 10, 64)
	if err != nil {
		return err
	}
	cnf.NetworkId = networkId

	a.base.SetConfig(cnf)
	//Modify the encryption algorithm to use the encryption algorithm on the chain
	if cnf.Cryptoer != ret.Cryptoer || cnf.Hasher != ret.Hasher {
		cnf.Cryptoer = ret.Cryptoer
		cnf.Hasher = ret.Hasher
		a.base.SetConfig(cnf)
		err = a.base.Init()
		if err != nil {
			return err
		}
	}

	return nil
}

type loginForm struct {
	EcosystemID int64  `json:"ecosystem_id"`
	Expire      int64  `json:"expire"`
	PublicKey   string `json:"public_key"`
	KeyID       string `json:"key_id"`
	Signature   string `json:"signature"`
	RoleID      int64  `json:"role_id"`
}

// Login
// roleId : Get the role id through the keyinfo interface
func (a *auth) Login(roleId int64) (err error) {
	var (
		sign []byte
		ret  loginResult
	)

	cnf := a.base.GetConfig()
	nonceSalt := fmt.Sprintf("LOGIN%d", cnf.NetworkId)
	sign, err = crypto.SignString(cnf.PrivateKey, nonceSalt+cnf.UID)
	if err != nil {
		return
	}

	form := &loginForm{}
	form.PublicKey = hex.EncodeToString(cnf.PublicKey)
	form.Signature = hex.EncodeToString(sign)
	form.EcosystemID = cnf.Ecosystem
	if roleId != 0 {
		form.RoleID = roleId
	}

	var req request.Request
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "login",
		Params:    []any{form},
	}
	req, err = a.base.NewMessage(message)
	if err != nil {
		return
	}

	err = a.base.GET(req, &ret)
	if err == nil {
		cnf.Token = ret.Token
		a.base.SetConfig(cnf)
	}
	return
}

func (a *auth) AutoLogin() (err error) {
	cnf := a.base.GetConfig()
	if cnf.Token != "" {
		if time.Unix(cnf.TokenExpireTime, 0).Sub(time.Now()) < time.Minute*10 {
			cnf.Token = ""
			a.base.SetConfig(cnf)
		} else {
			return nil
		}
	}
	err = a.GetUid()
	if err == nil {
		//default 0
		err = a.Login(0)
		if err == nil {
			err = a.GetUid()
		}
	}
	return
}

func (a *auth) GetAuthStatus() (*response.AuthStatusResponse, error) {
	var result response.AuthStatusResponse

	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getAuthStatus",
	}
	req, err := a.base.NewMessage(message)
	if err != nil {
		return &result, err
	}

	err = a.base.GET(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}
