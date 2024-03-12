package rest

import (
	"encoding/json"
	"net/http"
)

type apiError struct {
	error
	status  int
	message string
}

func (e apiError) Unwrap() error { return e.error }
func (e apiError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Status  int    `json:"status"`
		Message string `json:"message,omitempty"`
	}{e.status, e.message})
}

func newAPIError(err error, status int, msg string) error {
	if status == 0 {
		status = http.StatusInternalServerError
	}
	return apiError{
		error:   err,
		status:  status,
		message: msg,
	}
}
