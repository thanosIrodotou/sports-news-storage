package api

import (
	"context"
	"encoding/json"
	"net/http"

	ut "github.com/go-playground/universal-translator"
)

// Responder interface
type Responder interface {
	Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error
	RespondError(ctx context.Context, w http.ResponseWriter, err error) error
}

type JSONResponder struct {
	serviceName string
	translator  ut.Translator
}

// NewJSONResponder creates a new JSONResponder
func NewJSONResponder(sn string, tr ut.Translator) Responder {
	return &JSONResponder{
		serviceName: sn,
		translator:  tr,
	}
}

// Respond converts a Go value to JSON and serves it as a response back to the caller
func (r *JSONResponder) Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {
	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
	}

	return json.NewEncoder(w).Encode(&data)
}

// RespondError creates an error struct and serves it as a response back to the client. Extra data can be passed for more detailed logging of errors
func (r *JSONResponder) RespondError(ctx context.Context, w http.ResponseWriter, err error) error {
	statusCode := http.StatusInternalServerError
	errResp := ErrorResponse{
		Message: err.Error(),
		Details: err.Error(),
		Service: r.serviceName,
	}

	if err, ok := err.(Error); ok {
		errResp.Message = err.Error()
		errResp.Details = err.Error()
		errResp.Type = err.Code
		statusCode = err.statusCode
	}

	return r.Respond(ctx, w, errResp, statusCode)
}
