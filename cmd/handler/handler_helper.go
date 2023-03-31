package handler

import (
	"encoding/json"
)

type RareCircleError struct {
	Error string `json:"error"`
}

func InitializeError(errorMessage string) []byte {
	dapperError := RareCircleError{
		Error: errorMessage,
	}
	jsonResp, err := json.Marshal(dapperError)
	if err != nil {
		return []byte(err.Error())
	}

	return jsonResp
}
