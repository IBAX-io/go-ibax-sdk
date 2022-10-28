package response

type SendTxResult struct {
	Hashes map[string]string `json:"hashes"`
}

type TxStatusResult struct {
	BlockId int64 `json:"block_id"`
	Penalty int64 `json:"penalty"`
	Err     error `json:"err"`
}
