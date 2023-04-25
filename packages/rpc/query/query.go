package query

import (
	"errors"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/packages/modus"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
	"github.com/IBAX-io/go-ibax-sdk/packages/response"
	"strconv"
)

type query struct {
	modus.Base
}

func New(b modus.Base) modus.Query {
	return &query{Base: b}
}

func (q *query) GetKeyInfo(account string) (*response.KeyInfoResult, error) {
	var ret response.KeyInfoResult
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getKeyInfo",
		Params:    []any{account},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &ret, err
	}

	err = q.GET(req, &ret)
	if err != nil {
		return &ret, err
	}
	return &ret, nil
}

func (q *query) EcosystemInfo(ecosystem int64) (*response.EcosystemInfo, error) {
	var result response.EcosystemInfo
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "ecosystemInfo",
		Params:    []any{ecosystem},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &result, err
	}
	err = q.GET(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) EcosystemCount() (int64, error) {
	var result int64
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getEcosystemCount",
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return result, err
	}
	err = q.GET(req, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// SystemParams Find the specified name system params
// names Use‘,’ split multiple query parameters
// names if is zero value: get all
func (q *query) SystemParams(names string, params ...int) (*response.ParamsResult, error) {
	var result response.ParamsResult
	var offset, limit int
	if len(params) > 2 {
		return nil, errors.New("params invalid")
	}
	for k, v := range params {
		if k == 0 {
			offset = v
		} else {
			limit = v
		}
	}
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "systemParams",
		Params:    []any{names, request.Omitempty(offset), request.Omitempty(limit)},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &result, err
	}

	err = q.GET(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

// EcosystemParams Find the specified name ecosystem params
// names Use‘,’ split multiple query parameters
// names if is zero value: get all
func (q *query) EcosystemParams(ecosystem int64, names string, params ...int) (*response.ParamsResult, error) {
	var result response.ParamsResult
	var offset, limit int
	if len(params) > 2 {
		return nil, errors.New("params invalid")
	}
	for k, v := range params {
		if k == 0 {
			offset = v
		} else {
			limit = v
		}
	}
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getEcosystemParams",
		Params: []any{ecosystem, names,
			request.Omitempty(offset), request.Omitempty(limit)},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &result, err
	}

	err = q.GET(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) GetHistory(tableName string, tableId uint64) (*response.HistoryResult, error) {
	var result response.HistoryResult
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "history",
		Params:    []any{tableName, tableId},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &result, err
	}

	err = q.GET(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}
func (q *query) Balance(wallet string, ecosystemId int64) (*response.TokenBalanceResult, error) {
	var result response.TokenBalanceResult
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getBalance",
		Params:    []any{wallet, request.Omitempty(ecosystemId)},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &result, err
	}
	err = q.GET(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) GetMaxBlockID() (int64, error) {
	var result int64
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "maxBlockId",
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return result, err
	}
	err = q.GET(req, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (q *query) GetBlockInfo(id int64) (*response.BlockInfoResult, error) {
	var result response.BlockInfoResult
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getBlockInfo",
		Params:    []any{id},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &result, err
	}
	err = q.GET(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) AppParams(appid int64, names string, ecosystemId int64, params ...int) (*response.AppParamsResult, error) {
	type appParamsResult struct {
		App  int64                  `json:"app_id"`
		List []response.ParamResult `json:"list"`
	}

	var rlt appParamsResult
	var rets response.AppParamsResult
	var offset, limit int
	if len(params) > 2 {
		return nil, errors.New("params invalid")
	}
	for k, v := range params {
		if k == 0 {
			offset = v
		} else {
			limit = v
		}
	}
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "appParams",
		Params: []any{appid, request.Omitempty(ecosystemId),
			request.Omitempty(names), request.Omitempty(offset), request.Omitempty(limit)},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &rets, err
	}
	err = q.GET(req, &rlt)
	if err != nil {
		return &rets, err
	}
	rets.App = strconv.FormatInt(rlt.App, 10)
	rets.List = rlt.List
	return &rets, nil
}

// GetRow return where table id = id AND Login ecosystem = ecosystem, rowsName optional
func (q *query) GetRow(tableName string, id int64, columns string, whereColumn string) (*response.RowResult, error) {
	var result response.RowResult
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getRow",
		Params: []any{tableName, id, request.Omitempty(columns),
			request.Omitempty(whereColumn)},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &result, err
	}

	err = q.GET(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

// GetPageRow return page name = name and ecosystem id = login ecosystem id
func (q *query) GetPageRow(name string) (*response.PageResult, error) {
	var result response.PageResult
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getPageRow",
		Params:    []any{name},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &result, err
	}
	err = q.GET(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

// GetMenuRow return menu name = name and ecosystem id = login ecosystem id
func (q *query) GetMenuRow(name string) (*response.MenuResult, error) {
	var result response.MenuResult
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getMenuRow",
		Params:    []any{name},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &result, err
	}
	err = q.GET(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

// GetSnippetRow return snippet name = name and ecosystem id = login ecosystem id
func (q *query) GetSnippetRow(name string) (*response.SnippetResult, error) {
	var result response.SnippetResult
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getSnippetRow",
		Params:    []any{name},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &result, err
	}
	err = q.GET(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

// GetAppContent Get Obtain application-related information (including page, snippet, menu)
func (q *query) GetAppContent(appId int64) (*response.AppContentResult, error) {
	var result response.AppContentResult
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getAppContent",
		Params:    []any{appId},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &result, err
	}
	err = q.GET(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) BlocksTxInfo(blockId, count int64) (*map[int64][]response.BlockTxInfo, error) {
	var result map[int64][]response.BlockTxInfo
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getBlocksTxInfo",
		Params:    []any{blockId, count},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &result, err
	}
	err = q.GET(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) DetailedBlocks(blockId, count int64) (*map[int64]response.BlockDetailedInfo, error) {
	var result map[int64]response.BlockDetailedInfo
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "detailedBlocks",
		Params:    []any{blockId, count},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &result, err
	}
	err = q.GET(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func blockIdOrBlockHashValidate(req any) error {
	switch req.(type) {
	case string:
		if req.(string) == "" {
			return errors.New("params can't not be empty")
		}
	default:
		if cp, ok := req.(request.BlockIdOrHash); !ok {
			err := fmt.Errorf("params is not BlockIdOrHash type")
			return err
		} else {
			if cp.Id <= 0 && cp.Hash == "" {
				return errors.New("params can't not be empty")
			}
			if cp.Id > 0 && cp.Hash != "" {
				return errors.New("block id or block hash must be only choose one")
			}
		}
	}
	return nil
}

func (q *query) DetailedBlock(params any) (*response.BlockDetailedInfo, error) {
	err := blockIdOrBlockHashValidate(params)
	if err != nil {
		return nil, err
	}

	var result response.BlockDetailedInfo
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "detailedBlock",
		Params:    []any{params},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &result, err
	}
	err = q.GET(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) TransactionsCount() (int64, error) {
	var result int64
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getTxCount",
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return result, err
	}

	err = q.GET(req, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (q *query) BlockTxCount(blockIdOrHash any) (int64, error) {
	var result int64
	err := blockIdOrBlockHashValidate(blockIdOrHash)
	if err != nil {
		return result, err
	}
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getTransactionCount",
		Params:    []any{blockIdOrHash},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return result, err
	}

	err = q.GET(req, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (q *query) KeysCount() (int64, error) {
	var result int64
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getKeysCount",
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return result, err
	}
	err = q.GET(req, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (q *query) HonorNodesCount() (int64, error) {
	var result int64
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "honorNodesCount",
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return result, err
	}
	err = q.GET(req, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (q *query) GetTableCount(offset, limit int) (*response.TablesResult, error) {
	var v response.TablesResult
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getTableCount",
		Params:    []any{request.Omitempty(offset), request.Omitempty(limit)},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &v, err
	}
	err = q.GET(req, &v)
	if err != nil {
		return &v, err
	}
	return &v, nil
}

func (q *query) GetTable(tableName string) (*response.TableResult, error) {
	var result response.TableResult
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getTable",
		Params:    []any{tableName},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &result, err
	}
	err = q.GET(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

type sectionsForm struct {
	request.PaginatorForm
	Lang string `schema:"lang"`
}

func (q *query) GetSections(language string, offset, limit int) (*response.ListResult, error) {
	var params sectionsForm
	params.Lang = language
	params.Offset = offset
	params.Limit = limit

	var result response.ListResult
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getSections",
		Params:    []any{params},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &result, err
	}
	err = q.GET(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) GetIBAXConfig(option string) (*string, error) {
	var result string
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getConfig",
		Params:    []any{option},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &result, err
	}
	err = q.GET(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) GetVersion() (*string, error) {
	var result string
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getVersion",
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &result, err
	}
	err = q.GET(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) GetList(params request.GetList) (*response.ListResult, error) {
	var result response.ListResult
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getList",
		Params:    []any{params},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &result, err
	}

	err = q.POST(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) GetMemberInfo(account string, ecosystem int64) (*response.MemberInfo, error) {
	var result response.MemberInfo
	message := request.RequestParams{
		Namespace: request.NamespaceIBAX,
		Name:      "getMember",
		Params:    []any{account, ecosystem},
	}
	req, err := q.NewMessage(message)
	if err != nil {
		return &result, err
	}

	err = q.POST(req, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) GetBalance(wallet string) (*response.TokenBalanceResult, error) {
	return nil, response.NotSupportError
}

func (q *query) BlocksCount() (int64, error) {
	return 0, response.NotSupportError
}

func (q *query) DataVerify(tableName string, id int64, column, hash string) (string, error) {
	return "", response.NotSupportError
}

func (q *query) BinaryVerify(id int64, hash string) (any, error) {
	return nil, response.NotSupportError
}

func (q *query) EcosystemName(ecosystem int64) (*response.EcosystemNameResult, error) {
	return nil, response.NotSupportError
}

func (q *query) AppParam(appid int64, name string, ecosystem int64) (*response.ParamResult, error) {
	return nil, response.NotSupportError
}

func (q *query) EcosystemParam(ecosystem int64, name string) (*response.ParamResult, error) {
	return nil, response.NotSupportError
}

func (q *query) GetRowExtend(tableName, columns string, value string, rowsName string) (*response.RowResult, error) {
	return nil, response.NotSupportError
}

func (q *query) GetListWhere(tableName string, where any, columns, order string, page, limit int) (*response.ListResult, error) {
	return nil, response.NotSupportError
}

func (q *query) GetNodeListWhere(tableName string, where any, columns, order string, page, limit int) (*response.ListResult, error) {
	return nil, response.NotSupportError
}
