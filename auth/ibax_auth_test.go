package auth

import (
	"github.com/IBAX-io/go-ibax-sdk/base"
	"github.com/IBAX-io/go-ibax-sdk/config"
	"github.com/IBAX-io/go-ibax-sdk/pkg/common/crypto"
	"os"
	"strconv"
	"testing"
	"time"
)

func initConfig() *config.IbaxConfig {
	cnf := &config.IbaxConfig{}
	cnf.Ecosystem = 1
	cnf.JwtPrefix = "Bearer "
	cnf.PrivateKey = os.Getenv("IBAX_Private_Key1")
	cnf.ApiPath = "/api/v2/"
	cnf.ApiAddress = "http://localhost:7179"
	cnf.Cryptoer = "ECC_Secp256k1"
	cnf.Hasher = "KECCAK256"

	return cnf
}

func initErrConfig() *config.IbaxConfig {
	cnf := &config.IbaxConfig{}
	cnf.Ecosystem = 1
	cnf.JwtPrefix = "Bearer "
	cnf.PrivateKey = "6be1951fe2bc8ff0398fce6527aff9965e893f0f7e5ad2ef346"

	cnf.ApiPath = "/api/v2/"
	cnf.ApiAddress = "http://localhost:7179"
	cnf.Cryptoer = "ECC_Secp256k1"
	cnf.Hasher = "KECCAK256"

	return cnf
}

func TestIBAX_GetUid(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	c := NewClient(cnf, b)
	err := c.GetUid()
	if err != nil {
		t.Errorf("geuid failed :%s", err.Error())
	}
}

func TestIBAX_ErrorAccount(t *testing.T) {
	cnf := initErrConfig()
	b := base.NewClient(cnf)
	c := NewClient(cnf, b)
	err := c.GetUid()
	if err != nil {
		//t.Errorf("geuid failed :%s", err.Error())
	}
}

func TestIBAX_Login(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	c := NewClient(cnf, b)
	err := c.GetUid()
	if err != nil {
		t.Errorf("geuid failed :%s", err.Error())
	}
	err = c.Login("0")
	if err != nil {
		t.Errorf("login failed :%s", err.Error())
	}
}

func TestIBAX_AutoLogin(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	c := NewClient(cnf, b)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
	}
}

func TestIBAX_Refresh(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	c := NewClient(cnf, b)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}

	cfg := b.GetConfig()
	//fmt.Println("TokenExpireTime:", cfg.TokenExpireTime.String())
	cfg.TokenExpireTime = time.Minute * 5 //set token expire
	//fmt.Println("TokenExpireTime:", cfg.TokenExpireTime.String())
	if cfg.TokenExpireTime < time.Minute*10 {
		cfg.Token = ""

		err := c.AutoLogin()
		if err != nil {
			t.Errorf("auto login failed: %s", err.Error())
			return
		}
		//fmt.Println("TokenExpireTime:", cfg.TokenExpireTime.String())
	}
}

func TestIBAX_GetUID(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	c := NewClient(cnf, b)
	_, er := c.GetUID()
	if er != nil {
		t.Errorf("get uid failed :%s", er.Error())
		return
	}
}

func TestIBAX_LOGIN(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	c := NewClient(cnf, b)

	rets, er := c.GetUID()
	if er != nil {
		t.Errorf("login get uid failed :%s", er.Error())
		return
	}
	cfg := b.GetConfig()

	cfg.Token = rets.Token
	cfg.UID = rets.UID
	networkId, err := strconv.ParseInt(rets.NetworkID, 10, 64)
	if err != nil {
		t.Errorf("login Parse Int failed :%s", er.Error())
		return
	}
	cfg.NetworkId = networkId
	crypto.InitAsymAlgo(rets.Cryptoer)
	crypto.InitHashAlgo(rets.Hasher)
	if cfg.Hasher != rets.Hasher || cfg.Cryptoer != rets.Cryptoer {
		err = b.Init()
		if err != nil {
			return
		}
	}
	result, err := c.LOGIN("0")
	if err != nil {
		t.Errorf("login failed :%s", err.Error())
		return
	}
	cnf.Token = result.Token
}
