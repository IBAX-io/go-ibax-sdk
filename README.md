# IBAX Golang SDK
[![Go Reference](https://pkg.go.dev/badge/github.com/IBAX-io/go-ibax-sdk.svg)](https://pkg.go.dev/github.com/IBAX-io/go-ibax-sdk)

## It includes the following core components:
* **query** Query implementation
* **tx** transaction implementation
* **auth** implementation of authentication
* **pkg** contains core cryptographic, smart contracts and other useful functions


* The ones starting with `interface` are interfaces for different functions, the ones starting with `ibax` are the implementation of the interface,
* `client.go` encapsulates the functions of various requests, such as `getuid`, `login` uses the functions in auth


## Test
### Configure the address and private key information in `example` directory `init_test.go`,`initFounderTest` is found account, `initNormalTest` is normal account, `errAccountTest` is wrong private key account
### Execute The Test
```go
go test -coverpkg=./... -coverprofile="coverage.data" -v ./...

```
If there are use cases that fail the test, you can view them in the console
### View Test Report
```go
go tool cover -html="coverage.data" -o coverage.html

```
