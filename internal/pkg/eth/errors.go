package eth

import "fmt"

type ErrDecoding struct {
	message string
}

func NewErrDecoding(message string, args ...interface{}) *ErrDecoding {
	return &ErrDecoding{fmt.Sprintf(message, args...)}
}

func (e *ErrDecoding) Error() string {
	return e.message
}
