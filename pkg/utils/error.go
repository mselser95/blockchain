package utils

import "errors"

var (
	// ErrNotImplemented is returned when a method is not implemented.
	ErrNotImplemented = errors.New("not implemented")

	// ErrAlreadyStarted is returned when a manager is already started.
	ErrAlreadyStarted = errors.New("already started")

	// ErrClientNotStarted is returned when a client is not started.
	ErrClientNotStarted = errors.New("client not started")
)
