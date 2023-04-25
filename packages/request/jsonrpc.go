package request

import (
	"encoding/json"
	"errors"
)

type Namespace string

const (
	NamespaceSeparator = "."
)

const (
	NamespaceRPC Namespace = "rpc"

	NamespaceIBAX  Namespace = "ibax"
	NamespaceAdmin Namespace = "admin"
)

const JsonRPCVersion = "2.0"

type Request struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	ID      ID     `json:"id"`
	Params  Params `json:"params"`
}

type RequestParams struct {
	Namespace Namespace
	Name      string
	Params    []any
}

type BatchRequestParams struct {
	RequestParams
	Result any //response type
}

type BatchRequest struct {
	Req    *Request
	Result any   //response result
	Err    error //response error
}

func (r Request) MarshalJSON() ([]byte, error) {
	r2 := struct {
		Method  string `json:"method"`
		Params  Params `json:"params,omitempty"`
		ID      *ID    `json:"id,omitempty"`
		JSONRPC string `json:"jsonrpc"`
	}{
		Method:  r.Method,
		Params:  r.Params,
		JSONRPC: "2.0",
	}
	r2.ID = &r.ID
	return json.Marshal(r2)
}

// UnmarshalJSON implements json.Unmarshaler.
func (r *Request) UnmarshalJSON(data []byte) error {
	var r2 struct {
		Method string           `json:"method"`
		Params *json.RawMessage `json:"params,omitempty"`
		Meta   *json.RawMessage `json:"meta,omitempty"`
		ID     *ID              `json:"id"`
	}

	// Detect if the "params" field is JSON "null" or just not present
	// by seeing if the field gets overwritten to nil.
	r2.Params = &json.RawMessage{}

	if err := json.Unmarshal(data, &r2); err != nil {
		return err
	}
	r.Method = r2.Method
	if r2.Params == nil {
		r.Params = nil
	} else if len(*r2.Params) == 0 {
		r.Params = nil
	} else {
		err := json.Unmarshal(*r2.Params, &r.Params)
		if err != nil {
			return err
		}
	}

	if r2.Method == "" {
		return errors.New("request is missing method")
	}

	if r2.ID == nil {
		return errors.New("request is missing ID")
	} else {
		r.ID = *r2.ID
	}

	r.JSONRPC = "2.0"
	return nil
}
