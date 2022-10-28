package tx

import (
	"github.com/IBAX-io/go-ibax-sdk/auth"
	"github.com/IBAX-io/go-ibax-sdk/base"
	"github.com/IBAX-io/go-ibax-sdk/config"
	"os"
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

func TestIBAX_TxStatus(t *testing.T) {
	// Use the hash value returned by TestIBAX_SendTx
	cnf := initConfig()
	b := base.NewClient(cnf)
	c := NewClient(b)
	a := auth.NewClient(cnf, b)
	err := a.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
	}
	hash := "c46d69a0f8c6a51fac7786428a5cc3f5b59aef3a1ea8eb06b0e746dc6da8ea0d"
	_, err = c.TxStatus(hash, 10, 500*time.Millisecond)
	if err != nil {
		t.Errorf("get tx status failed :%s", err.Error())
		return
	}
}

func TestIBAX_TxsStatus(t *testing.T) {
	// Use the hash value returned by TestIBAX_SendTx
	cnf := initConfig()
	b := base.NewClient(cnf)
	c := NewClient(b)
	a := auth.NewClient(cnf, b)
	err := a.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
	}
	hashs := []string{
		"4240ce5c6c363740d795beffc3c60e04d7cf00997a2b33fe91a4b008443a5b33",
		"c46d69a0f8c6a51fac7786428a5cc3f5b59aef3a1ea8eb06b0e746dc6da8ea0d",
	}
	_, err = c.TxsStatus(hashs, 500*time.Millisecond)
	if err != nil {
		t.Errorf("get txs status failed :%s", err.Error())
		return
	}
}

func TestIBAX_GetTxInfo(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	c := NewClient(b)
	_, err := c.GetTxInfo("c46d69a0f8c6a51fac7786428a5cc3f5b59aef3a1ea8eb06b0e746dc6da8ea0d", false)
	if err != nil {
		t.Errorf("get tx info failed :%s", err.Error())
		return
	}
}

func TestIBAX_GetTxInfoMulti(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	c := NewClient(b)
	hashList := []string{
		"c46d69a0f8c6a51fac7786428a5cc3f5b59aef3a1ea8eb06b0e746dc6da8ea0d",
		"2fe382bbd0073c54c468969efaac1f69b088d3065d2729f67f5e8dfb470fdef5",
	}
	_, err := c.GetTxInfoMulti(hashList, false)
	if err != nil {
		t.Errorf("get tx info multi failed :%s", err.Error())
		return
	}
}
