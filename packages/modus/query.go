package modus

import (
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
	"github.com/IBAX-io/go-ibax-sdk/packages/response"
)

type Query interface {
	// GetKeyInfo
	// @method get
	// get key info
	GetKeyInfo(account string) (*response.KeyInfoResult, error)
	// Balance
	// @method get
	Balance(wallet string, ecosystem int64) (*response.TokenBalanceResult, error)
	GetBalance(wallet string) (*response.TokenBalanceResult, error)

	GetVersion() (*string, error)
	GetIBAXConfig(option string) (*string, error)

	//metrics
	EcosystemCount() (int64, error)
	GetMaxBlockID() (int64, error)
	TransactionsCount() (int64, error)
	KeysCount() (int64, error)
	HonorNodesCount() (int64, error)
	BlocksCount() (int64, error)

	//data verify
	DataVerify(tableName string, id int64, column, hash string) (string, error)
	BinaryVerify(id int64, hash string) (any, error)

	//block
	DetailedBlocks(block, count int64) (*map[int64]response.BlockDetailedInfo, error)
	GetBlockInfo(id int64) (*response.BlockInfoResult, error)
	BlocksTxInfo(block, count int64) (*map[int64][]response.BlockTxInfo, error)

	//ecosystem
	GetTableCount(offset, limit int) (*response.TablesResult, error)
	GetTable(tableName string) (*response.TableResult, error)
	GetSections(language string, offset, limit int) (*response.ListResult, error)
	GetPageRow(name string) (*response.PageResult, error)
	GetMenuRow(name string) (*response.MenuResult, error)
	GetSnippetRow(name string) (*response.SnippetResult, error)
	GetAppContent(appId int64) (*response.AppContentResult, error)

	//ecosystem
	EcosystemName(ecosystem int64) (*response.EcosystemNameResult, error)
	AppParams(appid int64, names string, ecosystem int64, params ...int) (*response.AppParamsResult, error)
	AppParam(appid int64, name string, ecosystem int64) (*response.ParamResult, error)
	EcosystemParam(ecosystem int64, name string) (*response.ParamResult, error)
	EcosystemParams(ecosystem int64, names string, params ...int) (*response.ParamsResult, error)
	SystemParams(names string, params ...int) (*response.ParamsResult, error)
	GetRow(tableName string, id int64, columns string, whereColumn string) (*response.RowResult, error)
	GetRowExtend(tableName, columns string, value string, rowsName string) (*response.RowResult, error)
	GetHistory(table string, id uint64) (*response.HistoryResult, error)
	GetList(params request.GetList) (*response.ListResult, error)
	GetListWhere(tableName string, where any, columns, order string, page, limit int) (*response.ListResult, error)
	GetNodeListWhere(tableName string, where any, columns, order string, page, limit int) (*response.ListResult, error)

	BlockTxCount(blockIdOrBlockHash any) (int64, error)
	DetailedBlock(blockIdOrBlockHash any) (*response.BlockDetailedInfo, error)
	EcosystemInfo(ecosystem int64) (*response.EcosystemInfo, error)
	GetMemberInfo(account string, ecosystem int64) (*response.MemberInfo, error)
}
