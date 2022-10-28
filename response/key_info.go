package response

import "github.com/shopspring/decimal"

type roleInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type notifyInfo struct {
	RoleID string `json:"role_id"`
	Count  int64  `json:"count"`
}

type keyEcosystemInfo struct {
	Ecosystem     string       `json:"ecosystem"`
	Name          string       `json:"name"`
	Roles         []roleInfo   `json:"roles,omitempty"`
	Notifications []notifyInfo `json:"notifications,omitempty"`
}

type KeyInfoResult struct {
	Account    string              `json:"account"`
	Ecosystems []*keyEcosystemInfo `json:"ecosystems"`
}

type TokenBalanceResult struct {
	Amount      string `json:"amount"`
	Money       string `json:"money"`
	Total       string `json:"total"`
	Utxo        string `json:"utxo"`
	TokenSymbol string `json:"token_symbol"`
}

type walletHistory struct {
	ID           int64           `json:"id"`
	SenderID     int64           `json:"sender_id"`
	SenderAdd    string          `json:"sender_add"`
	RecipientID  int64           `json:"recipient_id"`
	RecipientAdd string          `json:"recipient_add"`
	Amount       decimal.Decimal `json:"amount"`
	Comment      string          `json:"comment"`
	BlockID      int64           `json:"block_id"`
	TxHash       string          `json:"tx_hash"`
	CreatedAt    int64           `json:"created_at"`
	Money        string          `json:"money"`
}

type WalletHistoryResult struct {
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
	Total int64           `json:"total"`
	List  []walletHistory `json:"list"`
}
