package evm

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mselser95/blockchain/pkg/signer"
	"github.com/mselser95/blockchain/pkg/utils"
)

// Manager manages interactions with the EVM-compatible blockchain.
type Manager struct {
	client *ethclient.Client
	signer signer.TransactionSigner
}

// Connect establishes a connection to an EVM-compatible blockchain.
func (m *Manager) Connect(url string) error {
	var err error
	m.client, err = ethclient.Dial(url)
	return err
}

// GetBalance retrieves the balance of the specified address for a given token.
func (m *Manager) GetBalance(address string, token *utils.Token) (*big.Int, error) {
	account := common.HexToAddress(address)

	if token.Type == utils.Native {
		// Return the native currency balance (e.g., ETH)
		return m.client.BalanceAt(context.Background(), account, nil)
	} else if token.Type == utils.ERC20 {
		// Handle ERC20 token balance retrieval
		panic("not implemented")
	}

	return nil, errors.New("unsupported token type")
}

// ReadCall performs a read-only call to a contract on the EVM blockchain.
func (m *Manager) ReadCall(tx *utils.Transaction) (interface{}, error) {
	// Implementation remains the same as before
	return nil, nil
}

// SendTransaction signs and sends a transaction to the EVM blockchain.
func (m *Manager) SendTransaction(tx *utils.Transaction) (string, error) {
	if m.signer == nil {
		return "", errors.New("transaction signer is not set")
	}

	// Sign the transaction using the signer
	signedTx, err := m.signer.SignTransaction(tx)
	if err != nil {
		return "", err
	}

	// Convert the signed transaction data to a types.Transaction
	txData, err := hex.DecodeString(signedTx.Payload["signedTx"].(string))
	if err != nil {
		return "", err
	}
	var rawTx types.Transaction
	if err := rawTx.UnmarshalBinary(txData); err != nil {
		return "", err
	}

	// Send the signed transaction
	err = m.client.SendTransaction(context.Background(), &rawTx)
	if err != nil {
		return "", err
	}

	return rawTx.Hash().Hex(), nil
}

// GetTransactionDetails retrieves the details of a transaction by its ID.
func (m *Manager) GetTransactionDetails(txID string) (*utils.TransactionDetails, error) {
	// Implementation remains the same as before
	return nil, nil
}

// Start starts the Manager. This could involve setting up necessary resources.
func (m *Manager) Start() error {
	// Implement any startup logic if needed
	return nil
}

// Stop stops the Manager and cleans up resources.
func (m *Manager) Stop() error {
	// Implement any shutdown logic if needed
	return nil
}
