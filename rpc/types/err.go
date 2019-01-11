package types

import "fmt"

// Error wraps RPC errors, which contain an error code in addition to the message.
type Error interface {
	Error() string  // returns the message
	ErrorCode() int // returns the code
}

// request is for an unknown service
type MethodNotFoundError struct {
	service string
	method  string
}

func (e *MethodNotFoundError) ErrorCode() int { return -32601 }

func (e *MethodNotFoundError) Error() string {
	return fmt.Sprintf("The method %s%s%s does not exist/is not available", e.service, e.method)
}

// received message isn't a valid request
type InvalidRequestError struct{ message string }

func (e *InvalidRequestError) ErrorCode() int { return -32600 }

func (e *InvalidRequestError) Error() string { return e.message }

// received message is invalid
type InvalidMessageError struct{ message string }

func (e *InvalidMessageError) ErrorCode() int { return -32700 }

func (e *InvalidMessageError) Error() string { return e.message }

// unable to decode supplied params, or an invalid number of parameters
type InvalidParamsError struct{ message string }

func (e *InvalidParamsError) ErrorCode() int { return -32602 }

func (e *InvalidParamsError) Error() string { return e.message }

// logic error, callback returned an error
type CallbackError struct{ message string }

func (e *CallbackError) ErrorCode() int { return -32000 }

func (e *CallbackError) Error() string { return e.message }

// issued when a request is received after the server is issued to stop.
type ShutdownError struct{}

func (e *ShutdownError) ErrorCode() int { return -32000 }

func (e *ShutdownError) Error() string { return "server is shutting down" }
