package example

import (
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/client"
	"github.com/IBAX-io/go-ibax-sdk/tx/contract"
	"net/url"
	"testing"
)

// Contract Account Token Send
func TestIBAX_ContractTokenSend(t *testing.T) {
	cnf := initFounderTest()
	c := client.NewClient(cnf)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	form := url.Values{"Recipient": {"0000-0000-0000-0000-0000"}, "Amount": {"1000000000000"}}
	result, err := c.AutoCallContract("TokensSend", &form, "")
	if err != nil {
		t.Errorf("contract token send failed :%s", err.Error())
		return
	}

	fmt.Println("result:", result)
}

func TestIBAX_NewEcosystem(t *testing.T) {
	cnf := initFounderTest()
	c := client.NewClient(cnf)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}

	ecosystemName := "my first ecosystem"
	//params example 1
	//form := url.Values{"Name": {ecosystemName}}

	//params example 2
	form := contract.ContractParams{
		"Name": ecosystemName,
	}

	result, err := c.AutoCallContract("NewEcosystem", &form, "")
	if err != nil {
		t.Errorf("new ecosystem failed :%s", err.Error())
		return
	}
	fmt.Println("result:", *result)
}

func TestIBAX_Bytes(t *testing.T) {
	cnf := initFounderTest()
	c := client.NewClient(cnf)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}

	form := contract.ContractParams{
		"value": []byte{0x12},
	}

	result, err := c.AutoCallContract("test1", &form, "")
	if err != nil {
		t.Errorf("auto call contract failed :%s", err.Error())
		return
	}
	fmt.Println("result:", *result)
}
