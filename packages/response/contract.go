package response

import "encoding/json"

type GetContractResult struct {
	ID         uint32          `json:"id"`
	StateID    uint32          `json:"state"`
	TableID    string          `json:"tableid"`
	WalletID   string          `json:"walletid"`
	TokenID    string          `json:"tokenid"`
	Address    string          `json:"address"`
	Fields     []contractField `json:"fields"`
	Name       string          `json:"name"`
	AppId      uint32          `json:"app_id"`
	Ecosystem  uint32          `json:"ecosystem"`
	Conditions string          `json:"conditions"`
}

type contractField struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Optional bool   `json:"optional"`
}

type TxInfo struct {
	BlockId      int64          `json:"block_id"`
	BlockHash    string         `json:"block_hash"`
	Address      string         `json:"address"`
	Ecosystem    int64          `json:"ecosystem"`
	Hash         string         `json:"hash"`
	Expedite     string         `json:"expedite"`
	ContractName string         `json:"contract_name"`
	Params       map[string]any `json:"params"`
	CreatedAt    int64          `json:"created_at"`
	Size         string         `json:"size"`
	Status       int64          `json:"status"` //0:success 1:penalty
}

type TxInfoResult struct {
	BlockID string  `json:"blockid"`
	Confirm int     `json:"confirm"`
	Data    *TxInfo `json:"data"`
}

type MultiTxInfoResult struct {
	Results map[string]*TxInfoResult `json:"results"`
}

type ContentResult struct {
	Menu       string          `json:"menu"`
	MenuTree   json.RawMessage `json:"menutree"`
	Title      string          `json:"title"`
	Tree       json.RawMessage `json:"tree"`
	NodesCount int64           `json:"nodesCount"`
}

type HashResult struct {
	Hash string `json:"hash"`
}
