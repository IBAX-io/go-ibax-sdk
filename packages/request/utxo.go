package request

type UtxoType int

const (
	TypeContractToUTXO UtxoType = iota //Contract Account transfer UTXO Account
	TypeUTXOToContract                 //UTXO Account transfer Contract Account
	TypeTransfer
)

const (
	DestinationUtxo    = "UTXO"
	DestinationAccount = "Account"

	BlackHoleAddr = "0000-0000-0000-0000-0000"
)
