package solana

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mselser95/blockchain/pkg/signer"
	"github.com/mselser95/blockchain/pkg/utils"
)

// Manager manages interactions with the Solana blockchain.
type Manager struct {
	client *rpc.Client
	signer signer.TransactionSigner
}

// Connect establishes a connection to the Solana blockchain.
func (m *Manager) Connect(url string) error {
	m.client = rpc.New(url)
	return nil
}

// GetBalance retrieves the balance of the specified address for a given token.
func (m *Manager) GetBalance(address string, token *utils.Token) (*big.Int, error) {
	if token.Type == utils.Native {
		// Retrieve native SOL balance
		balance, err := m.client.GetBalance(context.Background(), address)
		if err != nil {
			return nil, err
		}
		return big.NewInt(int64(balance)), nil
	} else if token.Type == utils.SPLToken {
		// Handle SPL Token balance retrieval
		panic("not implemented")
	}

	return nil, errors.New("unsupported token type")
}

// ReadCall performs a read-only call to an account on the Solana blockchain.
func (m *Manager) ReadCall(tx *utils.Transaction) (interface{}, error) {
	// Solana read-only calls typically involve querying account info
	if tx.Payload["method"] == "GetAccountInfo" {
		account := tx.To
		accountInfo, err := m.client.GetAccountInfo(context.Background(), account)
		if err != nil {
			return nil, err
		}
		return accountInfo, nil
	}
	// Handle other read-only operations specific to Solana
	return nil, errors.New("unsupported method")
}

// SendTransaction signs and sends a transaction to the Solana blockchain.
func (m *Manager) SendTransaction(tx *utils.Transaction) (string, error) {
	if m.signer == nil {
		return "", errors.New("transaction signer is not set")
	}

	// Sign the transaction using the signer
	signedTx, err := m.signer.SignTransaction(tx)
	if err != nil {
		return "", err
	}

	// Send the signed transaction
	txID, err := m.client.SendTransactionWithOpts(context.Background(), signedTx.Payload["signedTx"].([]byte), false, rpc.CommitmentFinalized)
	if err != nil {
		return "", err
	}

	return txID, nil
}

// GetTransactionDetails retrieves the details of a transaction by its ID.
func (m *Manager) GetTransactionDetails(txID string) (*utils.TransactionDetails, error) {
	tx, err := m.client.GetTransaction(context.Background(), txID)
	if err != nil {
		return nil, err
	}

	if tx == nil {
		return nil, errors.New("transaction not found")
	}

	// Convert logs to a generic format
	logs := make([]utils.Log, len(tx.Meta.LogMessages))
	for i, logMsg := range tx.Meta.LogMessages {
		logs[i] = utils.Log{
			Address:     "",         // Solana logs may not always have an associated address
			Topics:      []string{}, // Solana logs do not have topics like EVM logs
			Data:        []byte(logMsg),
			BlockNumber: uint64(tx.Slot),
			TxHash:      txID,
			Index:       uint(i),
		}
	}

	status := "confirmed"
	if tx.Meta.Err != nil {
		status = "failed"
	}

	// Create the TransactionDetails struct
	details := &utils.TransactionDetails{
		ID:          txID,
		Status:      status,
		BlockNumber: uint64(tx.Slot),
		Timestamp:   time.Unix(int64(tx.BlockTime), 0),
		From:        "",  // Solana doesn't have a direct 'From' field in transactions
		To:          "",  // Solana transactions may involve multiple instructions and accounts
		Amount:      nil, // Needs to be parsed from transaction instructions if needed
		Fee:         big.NewInt(int64(tx.Meta.Fee)),
		Logs:        logs,
		Events:      map[string]interface{}{}, // Populate with relevant events if needed
	}

	return details, nil
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
