package example

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/packages/client"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/common/crypto"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/converter"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
	"github.com/IBAX-io/go-ibax-sdk/packages/response"
	"github.com/bitly/go-simplejson"
	"strconv"
	"testing"
)

func TestQuery_EcosystemCount(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()
	count, er := c.EcosystemCount()
	if er != nil {
		t.Errorf("ecosystem count failed :%s", er.Error())
		return
	}
	fmt.Printf("count:%d\n", count)

}
func TestQuery_SystemParams(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()
	//result, er := c.SystemParams("")
	_, er := c.SystemParams("number_of_nodes,taxes_size")
	if er != nil {
		t.Errorf("system params failed :%s", er.Error())
		return
	}

}
func TestQuery_EcosystemParam(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.EcosystemParam(1, "founder_account")
	if er != nil {
		t.Errorf("ecosystem param failed :%s", er.Error())
		return
	}

}
func TestQuery_EcosystemParams(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()
	//_, er := c.EcosystemParams(1, "changing_page,changing_menu")
	_, er := c.EcosystemParams(1, "")
	if er != nil {
		t.Errorf("ecosystem params failed :%s", er.Error())
		return
	}

}
func TestQuery_GetHistory(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.GetHistory("contracts", 15)
	if er != nil {
		t.Errorf("get history failed :%s", er.Error())
		return
	}

}

func TestQuery_GetBalance(t *testing.T) {
	c := client.NewClient(cnf)
	_, er := c.GetBalance("0666-0819-7161-7879-5186")
	if er != nil {
		t.Errorf("get balance failed :%s", er.Error())
		return
	}
}

func TestQuery_GetMaxBlockID(t *testing.T) {
	c := client.NewClient(cnf)
	_, er := c.GetMaxBlockID()
	if er != nil {
		t.Errorf("get max block id failed :%s", er.Error())
		return
	}

}

func TestQuery_GetBlockInfo(t *testing.T) {
	c := client.NewClient(cnf)
	v, er := c.GetBlockInfo(2)
	if er != nil {
		t.Errorf("get block info failed :%s", er.Error())
		return
	}
	fmt.Printf("%+v\n", v)
}

func TestQuery_Balance(t *testing.T) {
	c := client.NewClient(cnf)
	_, er := c.Balance("0666-0819-7161-7879-5186", 1)
	if er != nil {
		t.Errorf("get ibax balance failed :%s", er.Error())
		return
	}

}
func TestQuery_AppParams(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.AppParams(1, "", 1)
	if er != nil {
		t.Errorf("app params failed :%s", er.Error())
		return
	}
}

func TestQuery_AppParam(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.AppParam(1, "role_governancer", 1)
	if er != nil {
		t.Errorf("app param failed :%s", er.Error())
		return
	}

}

func TestQuery_EcosystemName(t *testing.T) {
	c := client.NewClient(cnf)
	_, er := c.EcosystemName(2)
	if er != nil {
		t.Errorf("get ecosystem name failed :%s", er.Error())
		return
	}

}

func TestQuery_GetList(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()
	var req request.GetList
	req.Name = "keys"
	req.Columns = "id,amount,ecosystem"
	req.Offset = 0
	req.Limit = 10
	req.Order = map[string]any{"ecosystem": -1}

	//jst, err := simplejson.NewJson([]byte(`{}`))
	//if err != nil {
	//	t.Errorf("new json failed :%s", err.Error())
	//	return
	//}

	j1, err := simplejson.NewJson([]byte(`{}`))
	if err != nil {
		t.Errorf("new json j1 failed :%s", err.Error())
		return
	}
	//j2, err := simplejson.NewJson([]byte(`{}`))
	//if err != nil {
	//	t.Errorf("new json j2 failed :%s", err.Error())
	//	return
	//}

	j1.Set("id", map[string]any{"$eq": -110277540701013350})
	//j2.Set("account", map[string]any{"$eq": "0666-0819-7161-7879-5186"})

	//jst.Set("$or", []any{j1})

	whereStr, _ := json.Marshal(j1)
	req.Where = string(whereStr)

	fmt.Printf("where:%s\n", req.Where)

	v, er := c.GetList(req)
	if er != nil {
		t.Errorf("get list failed :%s", er.Error())
		return
	}
	fmt.Printf("%+v\n", *v)
}

func TestQuery_GetRow(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()
	//_, er := c.GetRow("keys", 5555, "")
	result, er := c.GetRow("keys", -5476304279945383650, "account", "")
	if er != nil {
		t.Errorf("get row failed :%s", er.Error())
		return
	}
	fmt.Printf("result:%+v\n", result)
}

func TestQuery_DataVerify(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()
	table := "pages"
	id := int64(87)
	rowsName := "value"
	rets, er := c.GetRow(table, id, rowsName, "")
	if er != nil {
		t.Errorf("get row failed :%s", er.Error())
		return
	}
	cfg := c.GetConfig()

	if v, ok := rets.Value[rowsName]; ok {
		hash := crypto.Hash([]byte(v))

		tableName := converter.ParseTable(table, cfg.Ecosystem)
		//fileName := ""
		fileName := "data"
		result, err := c.DataVerify(tableName, id, rowsName, hex.EncodeToString(hash), fileName)
		if err != nil {
			t.Errorf("data verify failed :%s", err.Error())
			return
		}
		fmt.Printf("result:%+v\n", result)
	}

}

func TestQuery_BinaryVerify(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()
	id := int64(13)
	var req request.GetList
	req.Name = "@1binaries"
	req.Columns = "data"
	req.Offset = 0
	req.Limit = 1

	j1, err := simplejson.NewJson([]byte(`{}`))
	if err != nil {
		t.Errorf("new json j1 failed :%s", err.Error())
		return
	}
	j1.Set("id", map[string]any{"$eq": id})
	whereStr, _ := json.Marshal(j1)
	//fmt.Printf("wherestr:%s\n", req.Where)
	var rets *response.ListResult
	if !cnf.EnableRpc {
		req.Where = j1
		rets, err = c.GetNodeListWhere(req)
	} else {
		req.Where = string(whereStr)
		rets, err = c.GetList(req)
	}
	if err != nil {
		t.Errorf("get list failed:%s", err.Error())
		return
	}

	if rets.Count > 0 {
		for _, list := range rets.List {
			if v, ok := list[req.Columns]; ok {
				idStr, ok := list["id"]
				if !ok {
					continue
				}
				tableId, _ := strconv.ParseInt(idStr, 10, 64)
				if tableId == id {
					data, err := hex.DecodeString(v)
					if err != nil {
						t.Errorf("data decode failed:%s", err.Error())
						return
					}
					hash := crypto.Hash(data)
					var result any
					//fileName := ""
					fileName := "./data"
					result, err = c.BinaryVerify(id, hex.EncodeToString(hash), fileName)
					if err != nil {
						t.Errorf("binary verify failed :%s", err.Error())
						return
					}
					fmt.Printf("result:%+v\n", result)

				}
			}
		}
	}

}

func TestQuery_GetRowExtend(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.GetRowExtend("app_params", "name", "role_governancer", "app_id,name")
	//_, er := c.GetRowExtend("app_params", "name", "role_governancer", "")
	if er != nil {
		t.Errorf("get row extend failed :%s", er.Error())
		return
	}
}

func TestQuery_GetPageRow(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.GetPageRow("default_page")
	if er != nil {
		t.Errorf("get page row failed :%s", er.Error())
		return
	}
}

func TestQuery_GetMenuRow(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.GetMenuRow("default_menu")
	if er != nil {
		t.Errorf("get menu row failed :%s", er.Error())
		return
	}
}

func TestQuery_GetSnippetRow(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.GetSnippetRow("export_info")
	if er != nil {
		t.Errorf("get snippet row failed :%s", er.Error())
		return
	}
}

func TestQuery_GetAppContent(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()
	v, er := c.GetAppContent(1)
	if er != nil {
		t.Errorf("get App Content failed :%s", er.Error())
		return
	}
	fmt.Printf("%+v\n", *v)
}

func TestQuery_BlockTxInfo(t *testing.T) {
	c := client.NewClient(cnf)
	v, er := c.BlocksTxInfo(87, 2)
	if er != nil {
		t.Errorf("block tx info failed :%s", er.Error())
		return
	}
	fmt.Printf("%+v\n", *v)
}

func TestQuery_DetailBlocks(t *testing.T) {
	c := client.NewClient(cnf)
	_, er := c.DetailedBlocks(87, 1)
	if er != nil {
		t.Errorf("detailed blocks failed :%s", er.Error())
		return
	}
}

func TestQuery_BlocksCount(t *testing.T) {
	c := client.NewClient(cnf)
	_, er := c.BlocksCount()
	if er != nil {
		t.Errorf("blocks count failed :%s", er.Error())
		return
	}
}

func TestQuery_TransactionsCount(t *testing.T) {
	c := client.NewClient(cnf)
	_, er := c.TransactionsCount()
	if er != nil {
		t.Errorf("transactions count failed :%s", er.Error())
		return
	}
}

func TestQuery_KeysCount(t *testing.T) {
	c := client.NewClient(cnf)
	_, er := c.KeysCount()
	if er != nil {
		t.Errorf("keys count failed :%s", er.Error())
		return
	}
}

func TestQuery_HonorNodesCount(t *testing.T) {
	c := client.NewClient(cnf)
	_, er := c.HonorNodesCount()
	if er != nil {
		t.Errorf("honor nodes count failed :%s", er.Error())
		return
	}
}

func TestQuery_GetTables(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()
	v, er := c.GetTableCount(0, 2)
	if er != nil {
		t.Errorf("get table count failed :%s", er.Error())
		return
	}
	fmt.Printf("%+v\n", *v)
}

func TestQuery_GetTable(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.GetTable("app_params")
	if er != nil {
		t.Errorf("get table failed :%s", er.Error())
		return
	}
}

func TestQuery_GetSections(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.GetSections("en", 0, 10)
	if er != nil {
		t.Errorf("get Sections failed :%s", er.Error())
		return
	}
}

func TestQuery_GetIBAXConfig(t *testing.T) {
	c := client.NewClient(cnf)
	_, er := c.GetIBAXConfig("centrifugo")
	if er != nil {
		t.Errorf("get Centrifugo failed :%s", er.Error())
		return
	}
}

func TestQuery_GetVersion(t *testing.T) {
	c := client.NewClient(cnf)
	_, er := c.GetVersion()
	if er != nil {
		t.Errorf("get version failed :%s", er.Error())
		return
	}
}

func TestQuery_GetListWhere1(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()

	jst, err := simplejson.NewJson([]byte(`{}`))
	if err != nil {
		t.Errorf("new json failed :%s", err.Error())
		return
	}

	j1, err := simplejson.NewJson([]byte(`{}`))
	if err != nil {
		t.Errorf("new json j1 failed :%s", err.Error())
		return
	}
	j2, err := simplejson.NewJson([]byte(`{}`))
	if err != nil {
		t.Errorf("new json j2 failed :%s", err.Error())
		return
	}

	j1.Set("key", map[string]any{"$eq": "avatar"})
	j2.Set("account", map[string]any{"$eq": "0666-0819-7161-7879-5186"})

	jst.Set("$or", []any{j1, j2})

	var req request.GetList
	req.Where = jst
	req.Name = "buffer_data"
	req.Columns = "ecosystem,key,value,account"
	req.Order = "id desc"
	req.Offset = 0
	req.Limit = 10

	_, er := c.GetListWhere(req)
	if er != nil {
		t.Errorf("get list where failed :%s", er.Error())
		return
	}

}

func TestQuery_GetListWhere2(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()

	jst, err := simplejson.NewJson([]byte(`{}`))
	if err != nil {
		t.Errorf("new json failed :%s", err.Error())
		return
	}

	//jst.Set("key", map[string]any{"$eq": "avatar"})
	jst.Set("id", map[string]any{"$eq": "-110277540701013350"})
	var req request.GetList
	req.Where = jst
	req.Name = "keys"
	req.Order = "id desc"
	req.Offset = 0
	req.Limit = 10

	v, er := c.GetListWhere(req)
	if er != nil {
		t.Errorf("get list where failed :%s", er.Error())
		return
	}
	fmt.Printf("%+v\n", *v)

}

func TestQuery_GetNodeListWhere(t *testing.T) {
	c := client.NewClient(cnf)
	c.AutoLogin()

	jst, err := simplejson.NewJson([]byte(`{}`))
	if err != nil {
		t.Errorf("new json failed :%s", err.Error())
		return
	}

	jst.Set("block_id", map[string]any{"$eq": 1152})

	var req request.GetList
	req.Name = "history"
	req.Where = jst
	req.Order = map[string]any{"id": -1} //-1: desc 1: asc
	req.Offset = 0
	req.Limit = 10
	rets, er := c.GetNodeListWhere(req)
	if er != nil {
		t.Errorf("get node list where failed :%s", er.Error())
		return
	}
	fmt.Printf("count:%d\n", rets.Count)
	for i := 0; i < len(rets.List); i++ {
		fmt.Println(rets.List[i])
	}

}

func TestQuery_GetKeyInfo(t *testing.T) {
	c := client.NewClient(cnf)
	cnf = c.GetConfig()

	fmt.Printf("account:%s\n", cnf.Account)
	_, err := c.GetKeyInfo(cnf.Account)
	if err != nil {
		t.Errorf("get key info failed:%s", err.Error())
		return
	}
}

func TestQuery_EcosystemInfo(t *testing.T) {
	c := client.NewClient(cnf)

	v, err := c.EcosystemInfo(2)
	if err != nil {
		t.Errorf("get ecosystem info failed:%s", err.Error())
		return
	}
	fmt.Printf("v:%+v\n", *v)
}

func TestQuery_GetMemberInfo(t *testing.T) {
	c := client.NewClient(cnf)
	cnf = c.GetConfig()

	v, err := c.GetMemberInfo("0666-7782-2939-7671-3160", 1)
	if err != nil {
		t.Errorf("get member info failed:%s", err.Error())
		return
	}
	fmt.Printf("v:%+v\n", *v)
}

func TestQuery_DetailedBlock(t *testing.T) {
	c := client.NewClient(cnf)
	cnf = c.GetConfig()

	req := request.BlockIdOrHash{Id: 3}
	v, err := c.DetailedBlock(req)
	if err != nil {
		t.Errorf("detailed block failed:%s", err.Error())
		return
	}
	fmt.Printf("v:%+v\n", *v)
}

func TestQuery_BlockTxCount(t *testing.T) {
	c := client.NewClient(cnf)
	cnf = c.GetConfig()

	count, err := c.BlockTxCount("2")
	if err != nil {
		t.Errorf("detailed block failed:%s", err.Error())
		return
	}
	fmt.Printf("block tx count:%+v\n", count)
}

func TestQuery_GetAvatar(t *testing.T) {
	c := client.NewClient(cnf)
	cnf = c.GetConfig()

	account := "0666-0819-7161-7879-5186"
	ecosystemId := int64(1)
	fileName := fmt.Sprintf("%d-%s"+".png", ecosystemId, account)
	result, err := c.GetAvatar(account, ecosystemId, fileName)
	if err != nil {
		t.Errorf("get acatar failed:%s", err.Error())
		return
	}
	fmt.Printf("result:%+v\n", result)
}
