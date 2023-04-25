package modus

import (
	"github.com/IBAX-io/go-ibax-sdk/config"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
	"net/url"
)

type Base interface {
	Init() (err error)
	GetConfig() config.Config
	SetConfig(config.Config)
	Version() string
	AmountValidator(amount string) error
	ExpediteValidator(expedite string) error

	// RESTful API
	SendGet(url string, form *url.Values, result any) error
	SendPost(url string, form *url.Values, result any) error
	SendMultipart(url string, files map[string][]byte, result any) error

	// JSON-RPC Request
	NewMessage(params request.RequestParams) (request.Request, error)
	NewBatchMessage(requestPrams []request.BatchRequestParams) ([]request.BatchRequest, error)
	GET(form any, result any) error
	POST(form any, result any) error
}
