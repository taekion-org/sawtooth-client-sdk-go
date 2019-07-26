package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// SawtoothClientTransportRestError represents an error returned by SawtoothClientTransportRest.
type SawtoothClientTransportRestError struct {
	// Method is the HTTP method used for the request.
	Method        string
	// StatusCode is the HTTP status code returned from the request.
	StatusCode    int
	// ErrorResponse is the actual error data returned from the request.
	ErrorResponse errorRestResponse
}

// errorRestResponse represents a REST API reply when an error occurs.
type errorRestResponse struct {
	Error	struct {
		Code		int		`json:"code"`
		Title		string	`json:"title"`
		Message		string	`json:"message"`
	} `json:"error"`

	Valid	bool
}

// NewSawtoothClientTransportRestError constructs a SawtoothClientTransportRestError from an http.Response.
func NewSawtoothClientTransportRestError(response *http.Response) *SawtoothClientTransportRestError {
	error := &SawtoothClientTransportRestError{Method: response.Request.Method, StatusCode: response.StatusCode}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		error.ErrorResponse.Valid = false
	}

	err = json.Unmarshal(responseData, &error.ErrorResponse)
	if err != nil {
		error.ErrorResponse.Valid = false
	}

	error.ErrorResponse.Valid = true

	return error
}

// Error implements the error interface for SawtoothClientTransportRestError.
func (self *SawtoothClientTransportRestError) Error() string {
	msg := fmt.Sprintf("Sawtooth REST API Error: method=%s, status=%d", self.Method, self.StatusCode)

	if self.ErrorResponse.Valid {
		msg += fmt.Sprintf(", code=%d, title=%s", self.ErrorResponse.Error.Code, self.ErrorResponse.Error.Title)
	}

	return msg
}
