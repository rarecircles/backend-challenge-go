package rpc

import "fmt"

type ErrResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *ErrResponse) Error() string {
	return fmt.Sprintf("rpc error (code %d): %s", e.Code, e.Message)
}
