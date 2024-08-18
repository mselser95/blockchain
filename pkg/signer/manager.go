package signer

import (
	"github.com/mselser95/blockchain/pkg/utils"
)

// TransactionSigner defines the interface for signing blockchain transactions.
//
//go:generate mockgen -destination=../../internal/mock/signer/mock_signer.go -package=mock_signer -source=manager.go
type TransactionSigner interface {
	// SignTransaction signs the provided transaction.
	SignTransaction(tx *utils.Transaction) (*utils.Transaction, error)
}
