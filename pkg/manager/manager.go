package manager

import (
	"context"
	"math/big"

	"github.com/mselser95/blockchain/pkg/utils"
)

// BlockchainManager defines the common interface for blockchain interactions.
type BlockchainManager interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	GetBalance(ctx context.Context, address utils.Address, token utils.Token) (*big.Int, error)
	ReadCall(ctx context.Context, tx utils.Transaction) (interface{}, error)
	SendTransaction(ctx context.Context, tx utils.Transaction) (string, error)
	GetTransactionDetails(ctx context.Context, txID string) (*utils.TransactionDetails, error)
}
