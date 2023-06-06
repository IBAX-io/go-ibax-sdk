package response

import (
	"encoding/json"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
)

var jsonNull = json.RawMessage("null")

type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      request.ID      `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *jsonError      `json:"error,omitempty"`
}

type BatchResponse []Response

// MarshalJSON implements json.Marshaler and adds the "jsonrpc":"2.0"
// property.
func (r Response) MarshalJSON() ([]byte, error) {

	if r.Error != nil {
		response := struct {
			JSONRPC string     `json:"jsonrpc"`
			ID      request.ID `json:"id"`
			Error   any        `json:"error,omitempty"`
		}{
			JSONRPC: request.JsonRPCVersion,
			ID:      r.ID,
			Error:   r.Error,
		}

		return json.Marshal(response)
	} else {
		response := struct {
			JSONRPC string     `json:"jsonrpc"`
			ID      request.ID `json:"id"`
			Result  any        `json:"result,omitempty"`
		}{
			JSONRPC: request.JsonRPCVersion,
			ID:      r.ID,
			Result:  r.Result,
		}

		if response.Result == nil {
			response.Result = jsonNull
		}

		return json.Marshal(response)
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (r *Response) UnmarshalJSON(data []byte) error {
	type tmpType Response
	if err := json.Unmarshal(data, (*tmpType)(r)); err != nil {
		return err
	}
	return nil
}
