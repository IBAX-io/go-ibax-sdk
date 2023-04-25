package example

import (
	"github.com/IBAX-io/go-ibax-sdk/config"
	"os"
)

//var cnf = initApiTest()

var cnf = initJsonRpcTest()

//var cnf = errAccountTest()

func initApiTest() config.Config {
	var cnf config.Config
	cnf.Ecosystem = 1
	cnf.JwtPrefix = "Bearer "
	cnf.PrivateKey = os.Getenv("IBAX_Private_Key1")
	cnf.ApiPath = "/api/v2/"
	cnf.ApiAddress = "http://localhost:7079"
	cnf.Cryptoer = "ECC_Secp256k1"
	cnf.Hasher = "KECCAK256"
	cnf.EnableRpc = false

	return cnf
}

func initJsonRpcTest() config.Config {
	var cnf config.Config
	cnf.Ecosystem = 1
	cnf.JwtPrefix = "Bearer "
	cnf.PrivateKey = os.Getenv("IBAX_Private_Key2")

	cnf.ApiPath = "/api/v2/"
	cnf.ApiAddress = "http://localhost:7079"
	cnf.Cryptoer = "ECC_Secp256k1"
	cnf.Hasher = "KECCAK256"
	cnf.EnableRpc = true

	return cnf
}

func errAccountTest() config.Config {
	var cnf config.Config
	cnf.Ecosystem = 1
	cnf.JwtPrefix = "Bearer "
	cnf.PrivateKey = "6be1951fe2bc8ff0398fce6527aff9965e893f0f7e5ad2ef346"

	cnf.ApiPath = "/api/v2/"
	cnf.ApiAddress = "http://localhost:7079"
	cnf.Cryptoer = "ECC_Secp256k1"
	cnf.Hasher = "KECCAK256"

	return cnf
}
