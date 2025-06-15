package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInternalServerError = &sentinelError{statusCode: http.StatusInternalServerError, message: "oops! something went wrong"}
	ErrNotFound            = &sentinelError{statusCode: http.StatusNotFound, message: "resource not found"}
	ErrBadRequest          = &sentinelError{statusCode: http.StatusBadRequest, message: "bad request"}
)

type APIError interface {
	APIError() (int, string)
	Causer
}

type Causer interface {
	Cause() error
}
type sentinelError struct {
	statusCode int
	message    string
}

func (s sentinelError) Cause() error {
	return s
}

func (s sentinelError) Error() string {
	return s.message
}

func (s sentinelError) APIError() (int, string) {
	return s.statusCode, s.message
}

type sentinelWrappedError struct {
	err      error
	sentinel *sentinelError
}

func (s sentinelWrappedError) Error() string {
	return s.err.Error()
}

func (e sentinelWrappedError) Is(err error) bool {
	return e.sentinel == err
}

func (e sentinelWrappedError) APIError() (int, string) {
	code, message := e.sentinel.APIError()

	return code, message
}

func (e sentinelWrappedError) Unwrap() error {
	return e.err
}
func (e sentinelWrappedError) Cause() error {
	err := e.err
	for {
		wrappedErr := errors.Unwrap(err)
		if wrappedErr == nil {
			return err
		}
		err = wrappedErr
	}
}

func WrapError(err error, sentinel *sentinelError) error {
	if err == nil {
		return nil
	}

	if sentinel == nil {
		sentinel = ErrInternalServerError
	}

	var apiErr APIError
	if errors.As(err, &apiErr) {
		return fmt.Errorf("%w: %w", err, sentinel)
	}

	return sentinelWrappedError{err: err, sentinel: sentinel}

}

func WrapErrorWithCode(err error, code int, message string) error {
	return WrapError(err, &sentinelError{statusCode: code, message: message})
}

func Cause(err error) error {
	for {
		wrappedErr := errors.Unwrap(err)
		if wrappedErr == nil {
			return err
		}
		err = wrappedErr
	}
}
