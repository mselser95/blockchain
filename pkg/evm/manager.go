package evm

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strings"

	"github.com/mselser95/blockchain/pkg/manager"
	"github.com/mselser95/blockchain/pkg/signer"
	"github.com/mselser95/blockchain/pkg/utils"
)

// Manager manages interactions with the EVM-compatible blockchain.
type Manager struct {
	url           string
	client        ClientInterface
	signer        signer.TransactionSigner
	clientFactory ClientFactory
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
func (m *Manager) GetBalance(ctx context.Context, address utils.Address, token utils.Token) (*big.Int, error) {
	if m.client == nil {
		return nil, utils.ErrClientNotStarted
	}

	switch token.Type {
	case utils.Native:
		return m.getNativeBalance(ctx, address)
	case utils.ERC20:
		return m.getERC20Balance(ctx, address, token)
	default:
		return nil, fmt.Errorf("unsupported token type: %v", token.Type)
	}
}

func (m *Manager) getNativeBalance(ctx context.Context, address utils.Address) (*big.Int, error) {
	balance, err := m.client.BalanceAt(ctx, common.HexToAddress(address.String()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get native token balance: %w", err)
	}
	return balance, nil
}

func (m *Manager) getERC20Balance(ctx context.Context, address utils.Address, token utils.Token) (*big.Int, error) {
	if m.client == nil {
		return nil, utils.ErrClientNotStarted
	}

	// Parse the ABI
	parsedABI, err := abi.JSON(strings.NewReader(erc20Abi))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ERC20 ABI: %w", err)
	}

	// Prepare the input parameters
	tokenAddress := *token.Address
	contractAddress := common.HexToAddress(tokenAddress.String())
	tokenOwner := common.HexToAddress(address.String())

	// Pack the input parameters to match the "balanceOf" function in the ABI
	input, err := parsedABI.Pack("balanceOf", tokenOwner)
	if err != nil {
		return nil, fmt.Errorf("failed to pack parameters for balanceOf: %w", err)
	}

	// Call the contract
	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: input,
	}

	output, err := m.client.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call balanceOf: %w", err)
	}

	// Unpack the output
	var balance *big.Int
	err = parsedABI.UnpackIntoInterface(&balance, "balanceOf", output)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack balanceOf result: %w", err)
	}

	return balance, nil
}

// ReadCall performs a read-only call to a contract on the EVM blockchain.
func (m *Manager) ReadCall(ctx context.Context, tx utils.Transaction) (interface{}, error) {
	return nil, utils.ErrNotImplemented
}

// SendTransaction signs and sends a transaction to the EVM blockchain.
func (m *Manager) SendTransaction(ctx context.Context, tx utils.Transaction) (string, error) {
	return "", utils.ErrNotImplemented
}

// GetTransactionDetails retrieves the details of a transaction by its ID.
func (m *Manager) GetTransactionDetails(ctx context.Context, txID string) (*utils.TransactionDetails, error) {
	return nil, utils.ErrNotImplemented
}
