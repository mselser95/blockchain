package utils

import (
	"math/big"
	"time"
)

// TransactionType is an enumeration of supported transaction types.
type TransactionType int

const (
	// Transfer represents a simple value transfer transaction.
	Transfer TransactionType = iota
	// ContractCall represents a transaction that invokes a smart contract function.
	ContractCall
	// Stake represents a transaction that stakes tokens.
	Stake
	// Delegate represents a transaction that delegates tokens.
	Delegate
)

// TransactionStatus is an enumeration of possible transaction statuses.
type TransactionStatus int

const (
	// Pending represents a transaction that has been submitted but not yet confirmed.
	Pending TransactionStatus = iota
	// Confirmed represents a transaction that has been included in a block.
	Confirmed
	// Failed represents a transaction that has failed.
	Failed
)

// Transaction defines the interface for a blockchain transaction.
type Transaction interface {
	Hash() *TxHash
	From() Address
	To() Address
	Amount() *big.Int
	Type() *TransactionType
	Payload() map[string]interface{}

	Status() *TransactionStatus
	Timestamp() *time.Time
	BlockNumber() *uint64
	SignedTx() []byte

	// Validate checks if the transaction fields are valid.
	Validate() error

	// SetStatus updates the transaction status.
	SetStatus(status TransactionStatus)

	// SetBlockNumber updates the block number where the transaction was included.
	SetBlockNumber(blockNumber uint64)

	// SetSignedTx updates the signed transaction.
	SetSignedTx(signedTx []byte)

	// SetPayload updates the transaction payload.
	SetPayload(key string, value interface{})
}
