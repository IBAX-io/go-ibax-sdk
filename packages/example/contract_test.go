package example

import (
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/packages/client"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
	"github.com/IBAX-io/go-ibax-sdk/packages/response"
	"net/url"
	"testing"
)

// Contract Account Token Send
func TestIBAX_ContractTokenSend(t *testing.T) {
	c := client.NewClient(cnf)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	form := url.Values{"Recipient": {"0000-0000-0000-0000-0000"}, "Amount": {"1"}}
	result, err := c.AutoCallContract("@1TokensSend", &form, "")
	if err != nil {
		t.Errorf("contract token send failed :%s", err.Error())
		return
	}

	fmt.Printf("result:%+v\n", *result)
}

func TestIBAX_NewEcosystem(t *testing.T) {
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
	form := request.MapParams{
		"Name": ecosystemName,
	}

	result, err := c.AutoCallContract("NewEcosystem", &form, "")
	if err != nil {
		t.Errorf("new ecosystem failed :%s", err.Error())
		return
	}
	fmt.Println("result:", *result)
}

func TestIBAX_BytesParams(t *testing.T) {
	c := client.NewClient(cnf)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}

	form := request.MapParams{
		"value": []byte{0x12},
	}

	result, err := c.AutoCallContract("testContract", &form, "")
	if err != nil {
		t.Errorf("auto call contract failed :%s", err.Error())
		return
	}
	fmt.Println("result:", *result)
}

func TestIBAX_GetAuthStatus(t *testing.T) {
	c := client.NewClient(cnf)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}

	result, err := c.GetAuthStatus()
	if err != nil {
		t.Errorf("get auth status failed :%s", err.Error())
		return
	}
	fmt.Println("result:", *result)
}

func TestIBAX_GetContract(t *testing.T) {
	c := client.NewClient(cnf)
	err := c.AutoLogin()
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
	c := client.NewClient(cnf)
	err := c.AutoLogin()
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

func TestQuery_BatchRequest(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()

	var detailedBlock response.BlockDetailedInfo
	var maxBlockId string //error response struct
	batchParams := []request.BatchRequestParams{
		{
			RequestParams: request.RequestParams{
				Namespace: request.NamespaceIBAX,
				Name:      "getTransactionCount",
				Params:    []any{"50"},
			},
			//no response type specified
		},
		{
			RequestParams: request.RequestParams{
				Namespace: request.NamespaceIBAX,
				Name:      "detailedBlock",
				Params:    []any{"50"},
			},
			Result: &detailedBlock,
			//no response type specified
		},
		{
			RequestParams: request.RequestParams{
				Namespace: request.NamespaceIBAX,
				Name:      "maxBlockId",
			},
			Result: &maxBlockId,
			//Specifies the error response type
		},
	}
	batchRequest, err := c.NewBatchMessage(batchParams)
	if err != nil {
		t.Errorf("new batch message failed :%s", err.Error())
		return
	}
	err = c.POST(batchRequest, batchParams)
	if err != nil {
		t.Errorf("batch request failed:%s", err.Error())
		return
	}
	for _, v := range batchRequest {
		if v.Err != nil {
			fmt.Printf("req:%v,err:%s\n", v.Req, v.Err)
			continue
		}
		fmt.Printf("req:%+v,rlt:%+v\n", v.Req, v.Result)
	}

	fmt.Printf("detailedBlock:%+v\n", detailedBlock)
	fmt.Printf("maxBlockId:%+v\n", maxBlockId)

}
