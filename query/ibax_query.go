package query

import (
	"encoding/json"
	"errors"
	"fmt"
	response2 "github.com/IBAX-io/go-ibax-sdk/response"
	"net/url"
	"strconv"
)

func (chain *queryClient) GetKeyInfo(account string) (*response2.KeyInfoResult, error) {
	var ret response2.KeyInfoResult
	err := chain.SendGet(`keyinfo/`+account, nil, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// WalletHistory Find account transaction history
// search type (outcome:transfer, out income:transfer in, all:get all)
func (chain *queryClient) WalletHistory(searchType string, page, limit int64) (*response2.WalletHistoryResult, error) {
	var result response2.WalletHistoryResult
	form := url.Values{}
	reqUrl := fmt.Sprintf("walletHistory?searchType=%s&page=%d&limit=%d", searchType, page, limit)
	err := chain.SendGet(reqUrl, &form, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (chain *queryClient) EcosystemName(ecosystem int64) (*response2.EcosystemNameResult, error) {
	var result response2.EcosystemNameResult
	reqUrl := fmt.Sprintf("ecosystemname?id=%d", ecosystem)
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (chain *queryClient) EcosystemCount() (*response2.CountResult, error) {
	var result response2.CountResult
	form := url.Values{}
	ecosystemsUrl := fmt.Sprintf("/metrics/ecosystems")
	err := chain.SendGet(ecosystemsUrl, &form, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SystemParams Find the specified name system params
// names Use‘,’ split multiple query parameters
// names if is zero value: get all
func (chain *queryClient) SystemParams(names string) (*response2.ParamsResult, error) {
	var result response2.ParamsResult
	systemParamsUrl := fmt.Sprintf("systemparams?names=%s", names)
	err := chain.SendGet(systemParamsUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (chain *queryClient) EcosystemParam(ecosystem int64, name string) (*response2.ParamResult, error) {
	var result response2.ParamResult
	reqUrl := fmt.Sprintf("ecosystemparam/%s?ecosystem=%d", name, ecosystem)
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// EcosystemParams Find the specified name ecosystem params
// names Use‘,’ split multiple query parameters
// names if is zero value: get all
func (chain *queryClient) EcosystemParams(ecosystem int64, names string) (*response2.ParamsResult, error) {
	var result response2.ParamsResult
	ecosystemParamsUrl := fmt.Sprintf("ecosystemparams?ecosystem=%d&names=%s", ecosystem, names)
	err := chain.SendGet(ecosystemParamsUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
func (chain *queryClient) GetHistory(table string, id string) (*response2.HistoryResult, error) {
	var result response2.HistoryResult
	historyUrl := fmt.Sprintf("history/%s/%s", table, id)
	err := chain.SendGet(historyUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
func (chain *queryClient) GetBalance(wallet string) (*response2.TokenBalanceResult, error) {
	var result response2.TokenBalanceResult
	reqUrl := fmt.Sprintf("balance/%s", wallet)
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
func (chain *queryClient) GetMaxBlockID() (*response2.MaxBlockResult, error) {
	var result response2.MaxBlockResult
	maxBlockIdUrl := fmt.Sprintf("maxblockid")
	err := chain.SendGet(maxBlockIdUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (chain *queryClient) GetBlockInfo(id int64) (*response2.BlockInfoResult, error) {
	var rets response2.BlockInfoResult
	blockInfoUrl := fmt.Sprintf("block/%d", id)
	err := chain.SendGet(blockInfoUrl, nil, &rets)
	if err != nil {
		return nil, err
	}
	return &rets, nil
}
func (chain *queryClient) Balance(wallet string, ecosystem int64) (*response2.TokenBalanceResult, error) {
	var result response2.TokenBalanceResult
	balanceUrl := fmt.Sprintf("balance/%s?ecosystem=%d", wallet, ecosystem)
	err := chain.SendGet(balanceUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
func (chain *queryClient) AppParams(appid int64, names string, ecosystem int64) (*response2.AppParamsResult, error) {
	var rets response2.AppParamsResult
	appParamsUrl := fmt.Sprintf("appparams/%d?ecosystem=%d&names=%s", appid, ecosystem, names)
	err := chain.SendGet(appParamsUrl, nil, &rets)
	if err != nil {
		return nil, err
	}
	return &rets, nil
}
func (chain *queryClient) AppParam(appid int64, name string, ecosystem int64) (*response2.ParamResult, error) {
	var rets response2.ParamResult
	appParamUrl := fmt.Sprintf("appparam/%d/%s?ecosystem=%d", appid, name, ecosystem)
	err := chain.SendGet(appParamUrl, nil, &rets)
	if err != nil {
		return nil, err
	}
	return &rets, nil
}

// GetRow return where table id = id AND Login ecosystem = ecosystem, rowsName optional
func (chain *queryClient) GetRow(tableName string, id int64, rowsName string) (*response2.RowResult, error) {
	var result response2.RowResult
	getRowUrl := fmt.Sprintf("row/%s/%d", tableName, id)
	if rowsName != "" {
		getRowUrl += fmt.Sprintf("?columns=%s", rowsName)
	}
	err := chain.SendGet(getRowUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRowExtend return where columns = value AND Login ecosystem = ecosystem, rowsName optional
func (chain *queryClient) GetRowExtend(tableName, columns string, value string, rowsName string) (*response2.RowResult, error) {
	var result response2.RowResult
	getRowUrl := fmt.Sprintf("row/%s/%s/%s", tableName, columns, value)
	if rowsName != "" {
		getRowUrl += fmt.Sprintf("?columns=%s", rowsName)
	}
	err := chain.SendGet(getRowUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPageRow return page name = name and ecosystem id = login ecosystem id
func (chain *queryClient) GetPageRow(name string) (*response2.PageResult, error) {
	var result response2.PageResult
	getRowUrl := fmt.Sprintf("interface/page/%s", name)
	err := chain.SendGet(getRowUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMenuRow return menu name = name and ecosystem id = login ecosystem id
func (chain *queryClient) GetMenuRow(name string) (*response2.MenuResult, error) {
	var result response2.MenuResult
	getMenuUrl := fmt.Sprintf("interface/menu/%s", name)
	err := chain.SendGet(getMenuUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSnippetRow return snippet name = name and ecosystem id = login ecosystem id
func (chain *queryClient) GetSnippetRow(name string) (*response2.SnippetResult, error) {
	var result response2.SnippetResult
	getMenuUrl := fmt.Sprintf("interface/snippet/%s", name)
	err := chain.SendGet(getMenuUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (chain *queryClient) BlockTxInfo(count, blockId int64) (*map[int64][]response2.TxInfo, error) {
	var result map[int64][]response2.TxInfo
	reqUrl := fmt.Sprintf("blocks?count=%d&block_id=%d", count, blockId)
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (chain *queryClient) DetailedBlocks(count, blockId int64) (*map[int64]response2.BlockDetailedInfo, error) {
	var result map[int64]response2.BlockDetailedInfo
	reqUrl := fmt.Sprintf("detailed_blocks?count=%d&block_id=%d", count, blockId)
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (chain *queryClient) BlocksCount() (*response2.CountResult, error) {
	var result response2.CountResult
	reqUrl := fmt.Sprintf("metrics/blocks")
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (chain *queryClient) TransactionsCount() (*response2.CountResult, error) {
	var result response2.CountResult
	reqUrl := fmt.Sprintf("metrics/transactions")
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (chain *queryClient) EcosystemsCount() (*response2.CountResult, error) {
	var result response2.CountResult
	reqUrl := fmt.Sprintf("metrics/ecosystems")
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (chain *queryClient) KeysCount() (*response2.CountResult, error) {
	var result response2.CountResult
	reqUrl := fmt.Sprintf("metrics/keys")
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (chain *queryClient) HonorNodesCount() (*response2.CountResult, error) {
	var result response2.CountResult
	reqUrl := fmt.Sprintf("metrics/honornodes")
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (chain *queryClient) DataVerify(tableName string, id int64, column, hash string) (*string, error) {
	var result string
	reqUrl := fmt.Sprintf("data/%s/%d/%s/%s", tableName, id, column, hash)
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (chain *queryClient) BinaryVerify(id int64, hash string) (*any, error) {
	reqUrl := fmt.Sprintf("data/%d/data/%s", id, hash)
	var v any
	err := chain.SendGet(reqUrl, nil, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (chain *queryClient) GetTables(limit, offset int) (*response2.TablesResult, error) {
	var v response2.TablesResult
	reqUrl := fmt.Sprintf("tables?offset=%d&limit=%d", offset, limit)
	err := chain.SendGet(reqUrl, nil, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (chain *queryClient) GetTable(tableName string) (*response2.TableResult, error) {
	var v response2.TableResult
	reqUrl := fmt.Sprintf("table/%s", tableName)
	err := chain.SendGet(reqUrl, nil, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (chain *queryClient) GetList(tableName, columns string, limit, offset int64) (*response2.ListResult, error) {
	var result response2.ListResult
	getListUrl := fmt.Sprintf("list/%s?limit=%d&offset=%d&columns=%s", tableName, limit, offset, columns)
	err := chain.SendGet(getListUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (chain *queryClient) GetSections(language string, offset, limit int) (*response2.ListResult, error) {
	var result response2.ListResult
	reqUrl := fmt.Sprintf("sections?offset=%d&limit=%d", offset, limit)
	if language != "" {
		reqUrl += fmt.Sprintf("&lang=%s", language)
	}
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (chain *queryClient) GetCentrifugo() (*string, error) {
	var result string
	reqUrl := "config/centrifugo"
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (chain *queryClient) GetVersion() (*string, error) {
	var result string
	reqUrl := "version"
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (chain *queryClient) GetListWhere(tableName string, where any, columns, order string, page, limit int) (*response2.ListResult, error) {
	var result response2.ListResult
	reqUrl := fmt.Sprintf("listWhere/%s", tableName)

	form := &url.Values{}
	if order != "" {
		form.Set("order", order)
	}
	if columns != "" {
		form.Set("columns", columns)
	}
	if page < 0 || limit <= 0 {
		return nil, errors.New("params invalid")
	}
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	form.Set("offset", strconv.Itoa(offset))
	if where != nil {
		data, err := json.Marshal(where)
		if err != nil {
			return nil, err
		}
		form.Set("where", string(data))
	}

	err := chain.SendPost(reqUrl, form, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (chain *queryClient) GetNodeListWhere(tableName string, where any, columns, order string, page, limit int) (*response2.ListResult, error) {
	var result response2.ListResult
	reqUrl := fmt.Sprintf("nodelistWhere/%s", tableName)

	form := &url.Values{}
	if order != "" {
		form.Set("order", order)
	}
	if columns != "" {
		form.Set("columns", columns)
	}
	if page < 0 || limit <= 0 {
		return nil, errors.New("params invalid")
	}
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	form.Set("offset", strconv.Itoa(offset))
	if where != nil {
		data, err := json.Marshal(where)
		if err != nil {
			return nil, err
		}
		form.Set("where", string(data))
	}

	err := chain.SendPost(reqUrl, form, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
