package evm

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/mselser95/blockchain/pkg/manager"
	"github.com/mselser95/blockchain/pkg/signer"
	"github.com/mselser95/blockchain/pkg/utils"
)

// Manager manages interactions with the EVM-compatible blockchain.
type Manager struct {
	url           string
	client        ClientInterface
	signer        signer.TransactionSigner
	clientFactory ClientFactory // Use the factory interface
}

// NewManager creates a new Manager instance.
func NewManager(
	url string,
	signer signer.TransactionSigner,
	clientFactory ClientFactory,
) manager.BlockchainManager {
	return &Manager{
		url:           url,
		signer:        signer,
		clientFactory: clientFactory,
	}
}

// Start establishes a connection to the EVM-compatible blockchain.
func (m *Manager) Start(ctx context.Context) error {
	if m.client == nil {
		c, err := m.clientFactory.DialContext(ctx, m.url)
		if err != nil {
			// Check if the error is related to context cancellation
			if errors.Is(ctx.Err(), context.Canceled) {
				return fmt.Errorf("connection to EVM client at %s was canceled: %w", m.url, ctx.Err())
			}
			// Check if the error is related to context deadline exceeded
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return fmt.Errorf("connection to EVM client at %s timed out: %w", m.url, ctx.Err())
			}
			// For other types of errors, wrap them with additional context
			return fmt.Errorf("unable to connect to EVM client at %s: %w", m.url, err)
		}
		m.client = c
		return nil
	}
	return utils.ErrAlreadyStarted
}

// Stop stops the Manager and cleans up resources.
func (m *Manager) Stop(ctx context.Context) error {
	if m.client != nil {
		m.client.Close()
		return nil
	}
	return utils.ErrClientNotStarted
}

// GetBalance retrieves the balance of the specified address for a given token.
func (m *Manager) GetBalance(ctx context.Context, address string, token *utils.Token) (*big.Int, error) {
	return nil, utils.ErrNotImplemented
}

// ReadCall performs a read-only call to a contract on the EVM blockchain.
func (m *Manager) ReadCall(ctx context.Context, tx *utils.Transaction) (interface{}, error) {
	return nil, utils.ErrNotImplemented
}

// SendTransaction signs and sends a transaction to the EVM blockchain.
func (m *Manager) SendTransaction(ctx context.Context, tx *utils.Transaction) (string, error) {
	return "", utils.ErrNotImplemented
}

// GetTransactionDetails retrieves the details of a transaction by its ID.
func (m *Manager) GetTransactionDetails(ctx context.Context, txID string) (*utils.TransactionDetails, error) {
	return nil, utils.ErrNotImplemented
}
