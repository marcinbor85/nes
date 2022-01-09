package api

import (
	"fmt"
)

type RequestError struct {
	StatusCode int
	E          error
}

func (err *RequestError) Error() string {
	return fmt.Sprintf("StatusCode %d: %v", err.StatusCode, err.E)
}
