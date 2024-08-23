package utils

import (
	"fmt"
	"strings"
)

// IsError checks if the error message contains the specified error type.
func IsError(err error, errorType string) bool {
	return err != nil && strings.Contains(err.Error(), errorType)
}

// WrapError wraps an error with a predefined error message constant and additional context.
func WrapError(contextMsg string, baseErr ...error) error {
	if baseErr == nil {
		return fmt.Errorf("%s", contextMsg)
	}
	return fmt.Errorf("%s: %w", contextMsg, baseErr[0])
}

var (
	// ErrNotImplemented is returned when a method is not implemented.
	ErrNotImplemented = "not implemented"

	// ErrAlreadyStarted is returned when a manager is already started.
	ErrAlreadyStarted = "already started"

	// ErrClientNotStarted is returned when a client is not started.
	ErrClientNotStarted = "client not started"

	// ErrUnsupportedTokenType is returned when a token type is not supported.
	ErrUnsupportedTokenType = "unsupported token type"

	// EVM SPECIFIC ERRORS

	// ErrEVMInsufficientFunds is returned when an address is invalid.
	ErrEVMInsufficientFunds = "insufficient funds for gas * price + value"

	// ErrEVMMaxGasCapExceeded is returned when the maximum gas cap is exceeded.
	ErrEVMMaxGasCapExceeded = "max gas cap exceeded"

	// ErrEVMReplacementUnderpriced is returned when a replacement transaction is underpriced.
	ErrEVMReplacementUnderpriced = "replacement transaction underpriced"

	// ErrEVMNonceTooLow is returned when the nonce is too low.
	ErrEVMNonceTooLow = "nonce too low"

	// ErrEVMFailedToSendTransaction is returned when a transaction fails to send.
	ErrEVMFailedToSendTransaction = "failed to send transaction"

	// ErrEVMInvalidTransaction is returned when a transaction is invalid.
	ErrEVMInvalidTransaction = "invalid transaction"

	// ErrEVMFailedToSignTransaction is returned when a transaction fails to sign.
	ErrEVMFailedToSignTransaction = "failed to sign transaction"

	// ErrEVMInvalidPrivateKey is returned when a private key is invalid.
	ErrEVMInvalidPrivateKey = "invalid private key"

	// ErrEVMFailedToRetrieveTransaction is returned when a transaction fails to retrieve.
	ErrEVMFailedToRetrieveTransaction = "failed to retrieve transaction"

	// ErrEVMInvalidAddress is returned when an address is invalid.
	ErrEVMInvalidAddress = "invalid address"

	// ErrEVMInvalidHash is returned when a hash is invalid.
	ErrEVMInvalidHash = "invalid hash"
)
