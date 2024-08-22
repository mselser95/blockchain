package evm

import (
	"errors"
	"math/big"
	"time"

	"github.com/mselser95/blockchain/pkg/utils"
)

// BaseTransaction is a concrete implementation of the Transaction interface.
type BaseTransaction struct {
	TxHash        *utils.TxHash            // Transaction ID or hash
	FromAddress   utils.Address            // Sender address
	ToAddress     utils.Address            // Receiver address
	TxAmount      *big.Int                 // Amount to be transferred
	TxType        *utils.TransactionType   // Type of the transaction (e.g., transfer, contract call)
	TxStatus      *utils.TransactionStatus // Current status of the transaction
	TxTimestamp   *time.Time               // Timestamp of the transaction
	TxBlockNumber *uint64                  // Block number where the transaction was included
	TxPayload     map[string]interface{}   // Additional data specific to the transaction type or blockchain
	TxSignedTx    []byte                   // Signed transaction
}

// NewTransaction creates a new instance of BaseTransaction.
func NewTransaction(
	txHash *utils.TxHash,
	from utils.Address,
	to utils.Address,
	amount *big.Int,
	txType *utils.TransactionType,
	status *utils.TransactionStatus,
	timestamp *time.Time,
	blockNumber *uint64,
	gasLimit uint64,
	gasPrice *big.Int,
	chainId *big.Int,
	nonce uint64,
	data []byte,
) utils.Transaction {
	tx := &BaseTransaction{
		TxHash:        txHash,
		FromAddress:   from,
		ToAddress:     to,
		TxAmount:      amount,
		TxType:        txType,
		TxStatus:      status,
		TxTimestamp:   timestamp,
		TxBlockNumber: blockNumber,
		TxPayload:     make(map[string]interface{}),
	}

	// Only set non-zero and non-nil values to payload
	if gasLimit != 0 {
		tx.TxPayload["gasLimit"] = gasLimit
	}
	if gasPrice != nil {
		tx.TxPayload["gasPrice"] = gasPrice
	}
	if chainId != nil {
		tx.TxPayload["chainId"] = chainId
	}
	if nonce != 0 {
		tx.TxPayload["nonce"] = nonce
	}
	if data != nil {
		tx.TxPayload["data"] = data
	}

	return tx
}

// Hash returns the transaction ID or hash.
func (t *BaseTransaction) Hash() *utils.TxHash {
	return t.TxHash
}

// From returns the sender address.
func (t *BaseTransaction) From() utils.Address {
	return t.FromAddress
}

// To returns the receiver address.
func (t *BaseTransaction) To() utils.Address {
	return t.ToAddress
}

// Amount returns the amount to be transferred.
func (t *BaseTransaction) Amount() *big.Int {
	return t.TxAmount
}

// Type returns the type of the transaction.
func (t *BaseTransaction) Type() *utils.TransactionType {
	return t.TxType
}

// Status returns the current status of the transaction.
func (t *BaseTransaction) Status() *utils.TransactionStatus {
	return t.TxStatus
}

// Timestamp returns the timestamp of the transaction.
func (t *BaseTransaction) Timestamp() *time.Time {
	return t.TxTimestamp
}

// BlockNumber returns the block number where the transaction was included.
func (t *BaseTransaction) BlockNumber() *uint64 {
	return t.TxBlockNumber
}

// Payload returns additional data specific to the transaction type or blockchain.
func (t *BaseTransaction) Payload() map[string]interface{} {
	return t.TxPayload
}

// SignedTx returns the signed transaction bytes.
func (t *BaseTransaction) SignedTx() []byte {
	return t.TxSignedTx
}

// Validate checks if the transaction fields are valid.
func (t *BaseTransaction) Validate() error {
	if t.From() == nil || t.From().String() == "" {
		return errors.New("missing sender address")
	}

	if t.To() == nil || t.To().String() == "" {
		return errors.New("missing recipient address")
	}

	if t.Amount() == nil || t.Amount().Sign() <= 0 {
		return errors.New("invalid transaction amount")
	}

	// Validate payload fields if they exist
	if data, ok := t.TxPayload["data"]; ok {
		if _, ok := data.([]byte); !ok {
			return errors.New("payload 'data' must be a []byte")
		}
	}

	if gasPrice, ok := t.TxPayload["gasPrice"]; ok {
		if _, ok := gasPrice.(*big.Int); !ok {
			return errors.New("payload 'gasPrice' must be *big.Int")
		}
	} else {
		return errors.New("missing payload 'gasPrice'")
	}

	if gasLimit, ok := t.TxPayload["gasLimit"]; ok {
		if _, ok := gasLimit.(uint64); !ok {
			return errors.New("payload 'gasLimit' must be uint64")
		}
	} else {
		return errors.New("missing payload 'gasLimit'")
	}

	if nonce, ok := t.TxPayload["nonce"]; ok {
		if _, ok := nonce.(uint64); !ok {
			return errors.New("payload 'nonce' must be uint64")
		}
	} else {
		return errors.New("missing payload 'nonce'")
	}

	if chainId, ok := t.TxPayload["chainId"]; ok {
		if _, ok := chainId.(*big.Int); !ok {
			return errors.New("payload 'chainId' must be *big.Int")
		}
	} else {
		return errors.New("missing payload 'chainId'")
	}

	return nil
}

// SetStatus updates the transaction status.
func (t *BaseTransaction) SetStatus(status utils.TransactionStatus) {
	t.TxStatus = &status
}

// SetBlockNumber updates the block number where the transaction was included.
func (t *BaseTransaction) SetBlockNumber(blockNumber uint64) {
	t.TxBlockNumber = &blockNumber
}

// SetSignedTx updates the signed transaction.
func (t *BaseTransaction) SetSignedTx(signedTx []byte) {
	t.TxSignedTx = signedTx
}

// SetPayload updates the transaction payload.
func (t *BaseTransaction) SetPayload(key string, value interface{}) {
	t.TxPayload[key] = value
}
