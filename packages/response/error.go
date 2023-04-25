package response

import "errors"

var NotSupportError = errors.New("not support")

type jsonError struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data,omitempty"`
}

func (e *jsonError) Error() string {
	return e.Message
}
