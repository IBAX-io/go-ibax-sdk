# IBAX Golang SDK
[![Go Reference](https://pkg.go.dev/badge/github.com/IBAX-io/go-ibax-sdk.svg)](https://pkg.go.dev/github.com/IBAX-io/go-ibax-sdk)

## It includes the following core components:
- **config** package is the config

- **example** package is the test cases

- **client** package encapsulates the functions of various requests, such as `getuid`, `login` uses the functions in auth

- **modus** package is an interface

- **rpc** AND **api** package is the implementation of various interfaces, **rpc** is a JSON-RPC interface, 
**api** is a restful api interface, both of which use the http protocol.including Interfaces for query, authentication, 
and transaction sending, etc

- **pkg** contains core cryptographic, smart contracts and other useful functions

- **wallet** package is contains account creation, including HD wallet, mnemonic generation, private key generation, etc.


## Test
Configure the address and private key information in `init_test.go` in `example` directory, `initApiTest` is the 
configuration of `restful api` interface, `initJsonTest` is the configuration of `JSON-RPC` interface, `errAccountTest` 
is wrong account configuration


### Execute The Test
``` go
go test -v -cover -coverpkg="./..." "./packages/example" -coverprofile="coverage.data"

```

If there are use cases that fail the test, you can view them in the console
### View Test Report
``` go
go tool cover -html="coverage.data" -o coverage.html

```