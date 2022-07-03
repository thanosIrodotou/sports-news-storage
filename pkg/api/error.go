package api

import "net/http"

// ErrBadRequest represents an error message for bad requests
var ErrBadRequest = NewError(http.StatusText(http.StatusBadRequest), "errBadRequest", http.StatusBadRequest)

type Error struct {
	message    string
	Code       string
	statusCode int
}

func NewError(m, c string, s int) Error {
	return Error{
		message:    m,
		Code:       c,
		statusCode: s,
	}
}

func (e Error) Error() string {
	return e.message
}
