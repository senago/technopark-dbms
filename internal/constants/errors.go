package constants

import (
	"errors"
	"net/http"
)

// CodedError is an error wrapper which wraps errors with http status codes.
type CodedError struct {
	err  error
	code int
}

func (ce *CodedError) Error() string {
	return ce.err.Error()
}

func (ce *CodedError) Code() int {
	return ce.code
}

func NewCodedError(msg string, code int) *CodedError {
	return &CodedError{errors.New(msg), http.StatusNotFound}
}

var (
	// Not Found
	ErrDBNotFound = &CodedError{errors.New("Not found in db"), http.StatusNotFound}

	// User
	ErrUserAlreadyExists = &CodedError{errors.New("User with given nickname already exists"), http.StatusConflict}
)
