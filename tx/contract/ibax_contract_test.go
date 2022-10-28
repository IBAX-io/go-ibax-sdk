package contract

import (
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/auth"
	"github.com/IBAX-io/go-ibax-sdk/base"
	"github.com/IBAX-io/go-ibax-sdk/config"
	"github.com/IBAX-io/go-ibax-sdk/tx"
	"net/url"
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

func TestIBAX_GetContract(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	a := auth.NewClient(cnf, b)
	x := tx.NewClient(b)
	c := NewClient(cnf, b, x)
	err := a.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	result, er := c.GetContract("NewEcosystem")
	if er != nil {
		t.Errorf("get contract failed :%s", er.Error())
		return
	}
	fmt.Printf("result:%+v\n", result)

}

func TestIBAX_GetContracts(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	a := auth.NewClient(cnf, b)
	x := tx.NewClient(b)
	c := NewClient(cnf, b, x)
	err := a.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	_, er := c.GetContracts(1, 2)
	if er != nil {
		t.Errorf("get contract failed :%s", er.Error())
		return
	}

}

func TestIBAX_PrepareTx(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	a := auth.NewClient(cnf, b)
	x := tx.NewClient(b)
	c := NewClient(cnf, b, x)
	err := a.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	form := url.Values{"Recipient": {"0000-0000-0000-0000-0000"}, "Amount": {"1000"}}
	//form := url.Values{"MemberName": {"founder"}, "Amount": {"500"}}
	//form := url.Values{"MemberName": {"founder"}, "Amount": {"500"}, "Comment": {"Hello Dear"}}
	_, _, err = c.PrepareContractTx("TokensSend", &form)
	if err != nil {
		t.Errorf("prepare failed :%s", err.Error())
		return
	}
}

func TestIBAX_NewTransaction(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	a := auth.NewClient(cnf, b)
	x := tx.NewClient(b)
	c := NewClient(cnf, b, x)
	err := a.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	form := url.Values{"Recipient": {"0000-0000-0000-0000-0000"}, "Amount": {"1000"}}
	//form := url.Values{"MemberName": {"founder"}, "Amount": {"500"}}
	//form := url.Values{"MemberName": {"founder"}, "Amount": {"500"}, "Comment": {"Hello Dear"}}
	params, contractId, err := c.PrepareContractTx("TokensSend", &form)
	if err != nil {
		t.Errorf("prepare failed :%s", err.Error())
		return
	}
	_, _, err = c.NewContractTransaction(contractId, params, "")
	if err != nil {
		t.Errorf("Add Tx Data failed :%s", err.Error())
		return
	}
}

func TestIBAX_SendTx(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	a := auth.NewClient(cnf, b)
	x := tx.NewClient(b)
	c := NewClient(cnf, b, x)
	err := a.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	form := url.Values{"Recipient": {"0000-0000-0000-0000-0000"}, "Amount": {"1000"}}
	//form := url.Values{"MemberName": {"founder"}, "Amount": {"500"}}
	//form := url.Values{"MemberName": {"founder"}, "Amount": {"500"}, "Comment": {"Hello Dear"}}
	params, contractId, err := c.PrepareContractTx("TokensSend", &form)
	if err != nil {
		t.Errorf("prepare failed :%s", err.Error())
		return
	}
	data, hash, err := c.NewContractTransaction(contractId, params, "")
	if err != nil {
		t.Errorf("Add Tx Data failed :%s", err.Error())
		return
	}
	arrData := make(map[string][]byte)
	arrData[fmt.Sprintf("%x", hash)] = data

	_, err = x.SendTx(arrData)
	if err != nil {
		t.Errorf("Send Tx failed :%s", err.Error())
		return
	}
}

func TestIBAX_CallContract(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	a := auth.NewClient(cnf, b)
	x := tx.NewClient(b)
	c := NewClient(cnf, b, x)
	err := a.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	form := url.Values{"Recipient": {"0000-0000-0000-0000-0000"}, "Amount": {"500"}}
	params, contractId, err := c.PrepareContractTx("TokensSend", &form)
	if err != nil {
		t.Errorf("prepare form1 failed :%s", err.Error())
		return
	}
	data, hash, err := c.NewContractTransaction(contractId, params, "")
	if err != nil {
		t.Errorf("Add Tx Data failed :%s", err.Error())
		return
	}

	arrData := make(map[string][]byte)
	arrData[fmt.Sprintf("%x", hash)] = data

	hashMap, err := x.SendTx(arrData)
	if err != nil {
		t.Errorf("Send Tx failed :%s", err.Error())
		return
	}
	if hashMap != nil {
		var txHash string
		for _, v := range *hashMap {
			txHash = v
		}
		_, err := x.TxStatus(txHash, 10, time.Millisecond*500)
		if err != nil {
			t.Errorf("Tx Status failed:%s", err.Error())
			return
		}
	}

}

func TestIBAX_CallContracts(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	a := auth.NewClient(cnf, b)
	x := tx.NewClient(b)
	c := NewClient(cnf, b, x)
	err := a.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	form1 := url.Values{"Recipient": {"0000-0000-0000-0000-0000"}, "Amount": {"500"}}
	form2 := url.Values{"Recipient": {"1196-2490-5275-7101-3496"}, "Amount": {"200"}, "Comment": {"Hello Dear"}}
	params1, contractId1, err := c.PrepareContractTx("TokensSend", &form1)
	if err != nil {
		t.Errorf("prepare form1 failed :%s", err.Error())
		return
	}
	data1, hash1, err := c.NewContractTransaction(contractId1, params1, "")
	if err != nil {
		t.Errorf("Add Tx Data failed :%s", err.Error())
		return
	}

	params2, contractId2, err := c.PrepareContractTx("TokensSend", &form2)
	if err != nil {
		t.Errorf("prepare form2 failed :%s", err.Error())
		return
	}

	data2, hash2, err := c.NewContractTransaction(contractId2, params2, "1")
	if err != nil {
		t.Errorf("Add Tx Data failed :%s", err.Error())
		return
	}

	arrData := make(map[string][]byte)
	arrData[fmt.Sprintf("%x", hash1)] = data1
	arrData[fmt.Sprintf("%x", hash2)] = data2

	hashMap, err := x.SendTx(arrData)
	if err != nil {
		t.Errorf("Send Tx failed :%s", err.Error())
		return
	}
	if hashMap != nil {
		var hashs []string
		for _, v := range *hashMap {
			hashs = append(hashs, v)
		}
		_, err := x.TxsStatus(hashs, time.Millisecond*500)
		if err != nil {
			t.Errorf("Txs Status failed:%s", err.Error())
			return
		}
	}

}

func TestIBAX_AutoCallContract(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	a := auth.NewClient(cnf, b)
	x := tx.NewClient(b)
	c := NewClient(cnf, b, x)
	err := a.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	form := url.Values{"Recipient": {"0000-0000-0000-0000-0000"}, "Amount": {"1000000000000"}}
	_, err = c.AutoCallContract("TokensSend", &form, "")
	if err != nil {
		t.Errorf("auto call contract failed :%s", err.Error())
		return
	}

}
