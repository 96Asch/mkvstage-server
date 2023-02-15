package domain

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrType string

const (
	NotFound      ErrType = "ResourceNotFound"
	BadRequest    ErrType = "BadRequest"
	NotAuthorized ErrType = "NotAuthorized"
	Internal      ErrType = "Internal"
)

type Error struct {
	Type    ErrType `json:"type"`
	Message string  `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

// Status is a mapping errors to status codes
// Of course, this is somewhat redundant since
// our errors already map http status codes
func (e *Error) Status() int {
	switch e.Type {
	case NotAuthorized:
		return http.StatusUnauthorized
	case BadRequest:
		return http.StatusBadRequest
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// Status checks the runtime type
// of the error and returns an http
// status code if the error is model.Error
func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.Status()
	}
	return http.StatusInternalServerError
}

func NewObjectNotFoundErr(object string) *Error {
	return &Error{
		Type:    NotFound,
		Message: fmt.Sprintf("could not find resource with name: %v", object),
	}
}

func NewRecordNotFoundErr(key, val string) *Error {
	return &Error{
		Type:    NotFound,
		Message: fmt.Sprintf("could not find a record with %v = %v", key, val),
	}
}

func NewInternalErr() *Error {
	return &Error{
		Type:    Internal,
		Message: "something went wrong in the server",
	}
}

func NewNotAuthorizedErr(message string) *Error {
	return &Error{
		Type:    NotAuthorized,
		Message: message,
	}
}

func NewBadRequestErr(message string) *Error {
	return &Error{
		Type:    BadRequest,
		Message: message,
	}
}
