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

// Transaction represents an abstract transaction on a blockchain network.
type Transaction struct {
	// Transaction ID or hash
	ID string
	// Sender address
	From string
	// Receiver address
	To string
	// Amount to be transferred
	Amount *big.Int
	// Token identifier (e.g., token contract address or denomination)
	Token string
	// Type of the transaction (e.g., transfer, contract call)
	Type TransactionType
	// Current status of the transaction
	Status TransactionStatus
	// Transaction fee
	Fee *big.Int
	// Timestamp of the transaction
	Timestamp time.Time
	// Block number where the transaction was included
	BlockNumber uint64
	// Additional data specific to the transaction type or blockchain
	Payload map[string]interface{}
}
