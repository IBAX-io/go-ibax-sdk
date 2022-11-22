package example

import (
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/client"
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

type reqParams map[string]any

func (cp *reqParams) Get(key string) string {
	if _, ok := (*cp)[key]; !ok {
		return ""
	}
	return fmt.Sprintf("%v", (*cp)[key])
}

func (cp *reqParams) Set(key string, value any) {
	(*cp)[key] = value
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
	form := reqParams{
		"Name": ecosystemName,
	}

	result, err := c.AutoCallContract("NewEcosystem", &form, "")
	if err != nil {
		t.Errorf("new ecosystem failed :%s", err.Error())
		return
	}
	fmt.Println("result:", *result)
}
