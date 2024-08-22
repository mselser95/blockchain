package evm

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
	return utils.WrapError(utils.ErrAlreadyStarted)
}

// Stop stops the Manager and cleans up resources.
func (m *Manager) Stop(ctx context.Context) error {
	if m.client != nil {
		m.client.Close()
		return nil
	}
	return utils.WrapError(utils.ErrClientNotStarted)
}

// GetBalance retrieves the balance of the specified address for a given token.
func (m *Manager) GetBalance(ctx context.Context, address utils.Address, token utils.Token) (*big.Int, error) {
	if m.client == nil {
		return nil, utils.WrapError(utils.ErrClientNotStarted)
	}

	switch token.Type {
	case utils.Native:
		return m.getNativeBalance(ctx, address)
	case utils.ERC20:
		return m.getERC20Balance(ctx, address, token)
	default:
		errMsg := fmt.Sprintf("%s: %v", utils.ErrUnsupportedTokenType, token.Type)
		return nil, utils.WrapError(errMsg)
	}
}

// ReadCall performs a read-only call to a contract on the EVM blockchain.
func (m *Manager) ReadCall(ctx context.Context, tx utils.Transaction) (interface{}, error) {
	return nil, utils.WrapError(utils.ErrNotImplemented)
}

// EstimateGas estimates the gas required to execute a transaction on the EVM blockchain.
func (m *Manager) EstimateGas(ctx context.Context, tx utils.Transaction) (*big.Int, error) {
	return nil, utils.WrapError(utils.ErrNotImplemented)
}

// SendTransaction sends a transaction to the EVM blockchain.
func (m *Manager) SendTransaction(ctx context.Context, tx utils.Transaction) (string, error) {
	if m.client == nil {
		return "", utils.WrapError(utils.ErrClientNotStarted)
	}

	// Sign the transaction using the configured signer
	signedTx, err := m.signer.SignTransaction(tx)
	if err != nil {
		return "", utils.WrapError(utils.ErrEVMFailedToSignTransaction, err)
	}
	if signedTx == nil {
		return "", utils.WrapError(utils.ErrEVMInvalidTransaction)
	}

	// Deserialize the signed transaction from the transaction payload
	ethTx, ok := signedTx.Payload()["signedTransaction"].(*types.Transaction)
	if !ok || ethTx == nil {
		return "", utils.WrapError(utils.ErrEVMInvalidTransaction)
	}

	// Send the signed transaction to the Ethereum network
	err = m.client.SendTransaction(ctx, ethTx)
	if err != nil {
		// Handle known errors based on their messages
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "insufficient funds"):
			return "", utils.WrapError(utils.ErrEVMInsufficientFunds, err)
		case strings.Contains(errMsg, "exceeds block gas limit"):
			return "", utils.WrapError(utils.ErrEVMMaxGasCapExceeded, err)
		case strings.Contains(errMsg, "replacement transaction underpriced"):
			return "", utils.WrapError(utils.ErrEVMReplacementUnderpriced, err)
		case strings.Contains(errMsg, "nonce too low"):
			return "", utils.WrapError(utils.ErrEVMNonceTooLow, err)
		default:
			return "", utils.WrapError(utils.ErrEVMFailedToSendTransaction, err)
		}
	}

	// Return the transaction hash
	return ethTx.Hash().Hex(), nil
}

// GetTransactionDetails retrieves the details of a transaction by its ID.
func (m *Manager) GetTransactionDetails(ctx context.Context, txID string) (*utils.TransactionDetails, error) {
	return nil, utils.WrapError(utils.ErrNotImplemented, nil)
}

// Internal functions:
func (m *Manager) getNativeBalance(ctx context.Context, address utils.Address) (*big.Int, error) {
	balance, err := m.client.BalanceAt(ctx, common.HexToAddress(address.String()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get native token balance: %w", err)
	}
	return balance, nil
}

func (m *Manager) getERC20Balance(ctx context.Context, address utils.Address, token utils.Token) (*big.Int, error) {
	if m.client == nil {
		return nil, utils.WrapError(utils.ErrClientNotStarted, nil)
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
