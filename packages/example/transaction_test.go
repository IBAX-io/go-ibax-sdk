package example

import (
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/packages/client"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
	"net/url"
	"testing"
	"time"
)

func TestIBAX_NewTransaction(t *testing.T) {
	c := client.NewClient(cnf)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	form := url.Values{"Recipient": {"0000-0000-0000-0000-0000"}, "Amount": {"1000"}}
	//form := url.Values{"Recipient": {"0000-0000-0000-0000-0000"}, "Amount": {"500"}, "Comment": {"Test comment"}}
	params, contractId, err := c.PrepareContractTx("TokensSend", &form)
	if err != nil {
		t.Errorf("prepare failed :%s", err.Error())
		return
	}
	_, _, err = c.NewContractTransaction(contractId, params, "")
	if err != nil {
		t.Errorf("new contract transaction failed :%s", err.Error())
		return
	}
}

func TestIBAX_TxStatus(t *testing.T) {
	// Use the hash value returned by SendTx
	c := client.NewClient(cnf)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	hash := "354473e7e180915967e959fcb05a262baed192010a1eef05fed197cc782f4d25"

	rlts, err := c.TxStatus(hash, 10, 2*time.Second)
	if err != nil {
		t.Errorf("get tx status failed :%s", err.Error())
		return
	}
	fmt.Printf("v:%+v\n", rlts)
}

func TestIBAX_TxsStatus(t *testing.T) {
	// Use the hash value returned by SendTx
	c := client.NewClient(cnf)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	hashes := []string{
		//"94ae19cb22c68e8ad703b8e0f7db5fbfaad24c6295806e9a1f0fea7a1849c9b3",
		"354473e7e180915967e959fcb05a262baed192010a1eef05fed197cc782f4d25",
	}
	rlts, err := c.TxsStatus(hashes, 500*time.Millisecond)
	if err != nil {
		t.Errorf("get txs status failed :%s", err.Error())
		return
	}
	for _, v := range rlts {
		fmt.Printf("v:%+v\n", v)
	}
}

func TestIBAX_SendTx(t *testing.T) {
	c := client.NewClient(cnf)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	form := request.MapParams{
		"Recipient": "0000-0000-0000-0000-0000",
		"Amount":    "1000",
	}
	params, contractId, err := c.PrepareContractTx("TokensSend", &form)
	if err != nil {
		t.Errorf("prepare failed :%s", err.Error())
		return
	}
	data, hash, err := c.NewContractTransaction(contractId, params, "")
	if err != nil {
		t.Errorf("new contract transaction failed :%s", err.Error())
		return
	}
	arrData := make(map[string][]byte)
	arrData[fmt.Sprintf("%x", hash)] = data
	//arrData Support for multiple transactions

	_, err = c.SendTx(arrData)
	if err != nil {
		t.Errorf("Send Tx failed :%s", err.Error())
		return
	}
}

func TestIBAX_GetTxInfo(t *testing.T) {
	c := client.NewClient(cnf)
	_, err := c.GetTxInfo("c46d69a0f8c6a51fac7786428a5cc3f5b59aef3a1ea8eb06b0e746dc6da8ea0d", false)
	if err != nil {
		t.Errorf("get tx info failed :%s", err.Error())
		return
	}
}

func TestIBAX_GetTxInfoMulti(t *testing.T) {
	c := client.NewClient(cnf)
	hashList := []string{
		"c46d69a0f8c6a51fac7786428a5cc3f5b59aef3a1ea8eb06b0e746dc6da8ea0d",
		"19c6f68925df53d9ffac912f7e55bc8148f2f99fa8f075c0a12cfa78dc054225",
	}
	v, err := c.GetTxInfoMulti(hashList, false)
	if err != nil {
		t.Errorf("get tx info multi failed :%s", err.Error())
		return
	}
	fmt.Printf("v:%v\n", *v)
}
