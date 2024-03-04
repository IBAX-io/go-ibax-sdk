module github.com/IBAX-io/go-ibax-sdk

go 1.20

//https://github.com/btcsuite/btcd/btcutil
require (
	github.com/bitly/go-simplejson v0.5.0
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0
	github.com/ethereum/go-ethereum v1.10.25
	github.com/gogo/protobuf v1.3.2
	github.com/pkg/errors v0.9.1
	github.com/shopspring/decimal v1.3.1
	github.com/sirupsen/logrus v1.9.0
	github.com/tjfoc/gmsm v1.4.1
	github.com/tyler-smith/go-bip39 v1.0.1-0.20181017060643-dbb3b84ba2ef
	github.com/vmihailenco/msgpack/v5 v5.3.5
	golang.org/x/crypto v0.8.0
)

require github.com/btcsuite/btcd/btcutil v1.1.5

require (
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/btcsuite/btcd/chaincfg/chainhash v1.1.0 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
)

require (
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/btcsuite/btcd v0.23.5-0.20231215221805-96c9fd8078fd
	github.com/kr/pretty v0.3.1 // indirect
)
