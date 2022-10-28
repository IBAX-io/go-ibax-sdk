package response

type MaxBlockResult struct {
	MaxBlockID int64 `json:"max_block_id"`
}

type BlockInfoResult struct {
	Hash          []byte `json:"hash"`
	EcosystemID   int64  `json:"ecosystem_id"`
	KeyID         int64  `json:"key_id"`
	Time          int64  `json:"time"`
	Tx            int32  `json:"tx_count"`
	RollbacksHash []byte `json:"rollbacks_hash"`
	NodePosition  int64  `json:"node_position"`
}

type BlockHeaderInfo struct {
	BlockID      int64 `json:"block_id"`
	Time         int64 `json:"time"`
	KeyID        int64 `json:"key_id"`
	NodePosition int64 `json:"node_position"`
	Version      int   `json:"version"`
}

type TxDetailedInfo struct {
	Hash         []byte         `json:"hash"`
	ContractName string         `json:"contract_name"`
	Params       map[string]any `json:"params"`
	KeyID        int64          `json:"key_id"`
	Time         int64          `json:"time"`
	Type         byte           `json:"type"`
	Size         string         `json:"size"`
}

type TxInfo struct {
	Hash         []byte         `json:"hash"`
	ContractName string         `json:"contract_name"`
	Params       map[string]any `json:"params"`
	KeyID        int64          `json:"key_id"`
}

type BlockDetailedInfo struct {
	Header        BlockHeaderInfo  `json:"header"`
	Hash          []byte           `json:"hash"`
	NodePosition  int64            `json:"node_position"`
	KeyID         int64            `json:"key_id"`
	Time          int64            `json:"time"`
	TxCount       int32            `json:"tx_count"`
	Size          string           `json:"size"`
	RollbacksHash []byte           `json:"rollbacks_hash"`
	MerkleRoot    []byte           `json:"merkle_root"`
	BinData       []byte           `json:"bin_data"`
	StopCount     int              `json:"stop_count"`
	Transactions  []TxDetailedInfo `json:"transactions"`
}
