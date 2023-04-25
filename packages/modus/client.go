package modus

type Client interface {
	Base
	Authentication
	Contract
	Utxo
	Transaction
	Query
	Wallet
}
