package base

import (
	"github.com/IBAX-io/go-ibax-sdk/config"
	"log"
	"net/url"
)

type Base interface {
	Init() (err error)
	SendGet(url string, form *url.Values, v any) error
	SendPost(url string, form *url.Values, v any) error
	SendMultipart(url string, files map[string][]byte, v any) error
	GetConfig() *config.IbaxConfig
	SetConfig(cnf *config.IbaxConfig)
	Version() string
}

type baseClient struct {
	config *config.IbaxConfig
}

func NewClient(config *config.IbaxConfig) Base {
	b := &baseClient{config: config}
	if err := b.Init(); err != nil {
		log.Fatalf("new client init failed:%s\n", err.Error())
	}
	return b
}
