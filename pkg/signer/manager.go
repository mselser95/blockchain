package signer

import (
	"github.com/mselser95/blockchain/pkg/utils"
)

// TransactionSigner defines the interface for signing blockchain transactions.
type TransactionSigner interface {
	SignTransaction(tx *utils.Transaction) (*utils.Transaction, error)
}
