package response

import "encoding/json"

type GetContractResult struct {
	ID       uint32          `json:"id"`
	StateID  uint32          `json:"state"`
	TableID  string          `json:"tableid"`
	WalletID string          `json:"walletid"`
	TokenID  string          `json:"tokenid"`
	Address  string          `json:"address"`
	Fields   []contractField `json:"fields"`
	Name     string          `json:"name"`
}

type contractField struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Optional bool   `json:"optional"`
}

type txInfo struct {
	Block    string         `json:"block"`
	Contract string         `json:"contract"`
	Params   map[string]any `json:"params"`
}

type TxInfoResult struct {
	BlockID string  `json:"blockid"`
	Confirm int     `json:"confirm"`
	Data    *txInfo `json:"data"`
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
