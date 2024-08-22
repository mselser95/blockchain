package evm_test

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"testing"
	"time"

	"github.com/mselser95/blockchain/pkg/evm"

	"github.com/golang/mock/gomock"
	mock_evm "github.com/mselser95/blockchain/internal/mock/evm"
	mock_signer "github.com/mselser95/blockchain/internal/mock/signer"
	"github.com/mselser95/blockchain/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// To run all tests in this file from the root directory with coverage and verbosity:
// go test -cover ./pkg/evm

// To run this specific test from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm -run TestManager_Start_Success
func TestManager_Start_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_evm.NewMockClientInterface(ctrl)
	mockClientFactory := mock_evm.NewMockClientFactory(ctrl)
	mockSigner := mock_signer.NewMockTransactionSigner(ctrl)

	// Expect the factory to dial and return the mock client
	mockClientFactory.EXPECT().DialContext(gomock.Any(), "http://localhost:8545").Return(mockClient, nil)

	manager := evm.NewManager("http://localhost:8545", mockSigner, mockClientFactory)

	// Test the Start method with a successful connection
	err := manager.Start(context.Background())
	assert.NoError(t, err)
}

// To run this specific test from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm -run TestManager_Start_AlreadyStarted
func TestManager_Start_AlreadyStarted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClientFactory := mock_evm.NewMockClientFactory(ctrl)
	mockSigner := mock_signer.NewMockTransactionSigner(ctrl)
	mockClient := mock_evm.NewMockClientInterface(ctrl)
	mockClientFactory.EXPECT().DialContext(gomock.Any(), "http://localhost:8545").Return(mockClient, nil)

	manager := evm.NewManager("http://localhost:8545", mockSigner, mockClientFactory)

	// Test the Start method when the manager is already started
	err := manager.Start(context.Background())
	assert.NoError(t, err)
	err = manager.Start(context.Background())
	assert.ErrorContains(t, err, utils.ErrAlreadyStarted)
}

// To run this specific test from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm -run TestManager_Start_CancelledContext
func TestManager_Start_CancelledContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClientFactory := mock_evm.NewMockClientFactory(ctrl)
	mockSigner := mock_signer.NewMockTransactionSigner(ctrl)

	manager := evm.NewManager("http://localhost:8545", mockSigner, mockClientFactory)

	// Create a context that is already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Expect the factory not to be called since the context is cancelled
	mockClientFactory.EXPECT().DialContext(gomock.Any(), gomock.Any()).Return(nil, context.Canceled).Times(1)

	// Test the Start method with a cancelled context
	err := manager.Start(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), context.Canceled.Error())
}

// To run this specific test from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm -run TestManager_Start_DeadlineExceeded
func TestManager_Start_DeadlineExceeded(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClientFactory := mock_evm.NewMockClientFactory(ctrl)
	mockSigner := mock_signer.NewMockTransactionSigner(ctrl)

	manager := evm.NewManager("http://localhost:8545", mockSigner, mockClientFactory)

	// Create a context with a very short deadline
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	time.Sleep(2 * time.Millisecond)

	// Expect the factory to be called and simulate an error due to the deadline being exceeded
	mockClientFactory.EXPECT().DialContext(gomock.Any(), "http://localhost:8545").Return(nil, context.DeadlineExceeded)

	// Test the Start method with a deadline-exceeded context
	err := manager.Start(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), context.DeadlineExceeded.Error())
}

// To run this specific test from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm -run TestManager_Start_DialError
func TestManager_Start_DialError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClientFactory := mock_evm.NewMockClientFactory(ctrl)
	mockSigner := mock_signer.NewMockTransactionSigner(ctrl)

	manager := evm.NewManager("http://localhost:8545", mockSigner, mockClientFactory)

	// Simulate a network error during the connection attempt
	dialError := errors.New("network error")
	mockClientFactory.EXPECT().DialContext(gomock.Any(), "http://localhost:8545").Return(nil, dialError)

	// Test the Start method with a network error during dialing
	err := manager.Start(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to connect to EVM client")
	assert.Contains(t, err.Error(), "network error")
}

// To run this specific test from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm -run TestManager_Stop_Success
func TestManager_Stop_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_evm.NewMockClientInterface(ctrl)
	mockClientFactory := mock_evm.NewMockClientFactory(ctrl)
	mockSigner := mock_signer.NewMockTransactionSigner(ctrl)
	mockClientFactory.EXPECT().DialContext(gomock.Any(), "http://localhost:8545").Return(mockClient, nil)

	manager := evm.NewManager("http://localhost:8545", mockSigner, mockClientFactory)

	// Start the manager to initialize the client
	manager.Start(context.Background())

	// Expect the client to be closed when Stop is called
	mockClient.EXPECT().Close().Times(1)

	// Test the Stop method
	err := manager.Stop(context.Background())
	assert.NoError(t, err)
}

// To run this specific test from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm -run TestManager_Stop_NoClient
func TestManager_Stop_NoClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClientFactory := mock_evm.NewMockClientFactory(ctrl)
	mockSigner := mock_signer.NewMockTransactionSigner(ctrl)

	manager := evm.NewManager("http://localhost:8545", mockSigner, mockClientFactory)

	// Test the Stop method when there is no client (i.e., manager was never started)
	err := manager.Stop(context.Background())
	assert.ErrorContains(t, err, utils.ErrClientNotStarted)
}

// To run this specific test from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm -run TestManager_Stop_CanceledContext
func TestManager_Stop_CanceledContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_evm.NewMockClientInterface(ctrl)
	mockClientFactory := mock_evm.NewMockClientFactory(ctrl)
	mockSigner := mock_signer.NewMockTransactionSigner(ctrl)
	mockClientFactory.EXPECT().DialContext(gomock.Any(), "http://localhost:8545").Return(mockClient, nil)

	manager := evm.NewManager("http://localhost:8545", mockSigner, mockClientFactory)

	// Start the manager to initialize the client
	manager.Start(context.Background())

	// Expect the client to be closed when Stop is called
	mockClient.EXPECT().Close().Times(1)

	// Create a context that is already canceled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Test the Stop method with a canceled context
	err := manager.Stop(ctx)
	assert.NoError(t, err)
}

// TestManager_GetBalance_Native_Success tests the GetBalance method for a native token (e.g., ETH).
// go test -v -cover ./pkg/evm -run TestManager_GetBalance_Native_Success
func TestManager_GetBalance_Native_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_evm.NewMockClientInterface(ctrl)
	mockClientFactory := mock_evm.NewMockClientFactory(ctrl)
	mockSigner := mock_signer.NewMockTransactionSigner(ctrl)
	mockClientFactory.EXPECT().DialContext(gomock.Any(), "http://localhost:8545").Return(mockClient, nil)

	manager := evm.NewManager("http://localhost:8545", mockSigner, mockClientFactory)

	// Start the manager to initialize the client
	manager.Start(context.Background())

	address, err := evm.NewAddress("0x32Be343B94f860124dC4fEe278FDCBD38C102D88", utils.Ethereum)
	assert.NoError(t, err)

	mockClient.EXPECT().BalanceAt(gomock.Any(), common.HexToAddress(address.String()), nil).
		Return(big.NewInt(1000000000000000000), nil)

	token := utils.Token{Type: utils.Native, Symbol: "ETH"}

	balance, err := manager.GetBalance(context.Background(), address, token)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(1000000000000000000), balance)
}

// TestManager_GetBalance_ERC20_Success tests the GetBalance method for an ERC20 token.
// go test -v -cover ./pkg/evm -run TestManager_GetBalance_ERC20_Success
func TestManager_GetBalance_ERC20_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_evm.NewMockClientInterface(ctrl)
	mockClientFactory := mock_evm.NewMockClientFactory(ctrl)
	mockSigner := mock_signer.NewMockTransactionSigner(ctrl)
	mockClientFactory.EXPECT().DialContext(gomock.Any(), "http://localhost:8545").Return(mockClient, nil)

	balanceExpected := big.NewInt(1000000000000000000)

	// Convert to a 32-byte array (padded with zeros at the start)
	output := common.LeftPadBytes(balanceExpected.Bytes(), 32)

	mockClient.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return(output, nil)

	manager := evm.NewManager("http://localhost:8545", mockSigner, mockClientFactory)

	// Start the manager to initialize the client
	manager.Start(context.Background())

	address, err := evm.NewAddress("0x32Be343B94f860124dC4fEe278FDCBD38C102D88", utils.Ethereum)
	assert.NoError(t, err)

	tokenAddress, err := evm.NewAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F", utils.Ethereum)
	assert.NoError(t, err)

	token := utils.Token{Type: utils.ERC20, Address: &tokenAddress, Symbol: "DAI"}

	balance, err := manager.GetBalance(context.Background(), address, token)
	assert.NoError(t, err)
	assert.Equal(t, balanceExpected, balance)
}

// TestManager_GetBalance_UnsupportedToken tests the GetBalance method for an unsupported token type.
// go test -v -cover ./pkg/evm -run TestManager_GetBalance_UnsupportedToken
func TestManager_GetBalance_UnsupportedToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClientFactory := mock_evm.NewMockClientFactory(ctrl)
	mockSigner := mock_signer.NewMockTransactionSigner(ctrl)
	mockClient := mock_evm.NewMockClientInterface(ctrl)
	mockClientFactory.EXPECT().DialContext(gomock.Any(), "http://localhost:8545").Return(mockClient, nil)

	manager := evm.NewManager("http://localhost:8545", mockSigner, mockClientFactory)

	address, err := evm.NewAddress("0x32Be343B94f860124dC4fEe278FDCBD38C102D88", utils.Ethereum)
	assert.NoError(t, err)

	err = manager.Start(context.Background())
	assert.NoError(t, err)

	token := utils.Token{Type: utils.CosmosDenom, Symbol: "ATOM"}

	balance, err := manager.GetBalance(context.Background(), address, token)
	assert.Error(t, err)
	assert.Nil(t, balance)
	assert.ErrorContains(t, err, utils.ErrUnsupportedTokenType)
}

// TestMangaer_SendTransaction_ClientNotStarted tests the SendTransaction method when the client is not started.
// go test -v -cover ./pkg/evm -run TestManager_SendTransaction_ClientNotStarted
func TestManager_SendTransaction_ClientNotStarted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClientFactory := mock_evm.NewMockClientFactory(ctrl)
	mockSigner := mock_signer.NewMockTransactionSigner(ctrl)

	manager := evm.NewManager("http://localhost:8545", mockSigner, mockClientFactory)

	_, err := manager.SendTransaction(context.Background(), &evm.BaseTransaction{})
	assert.ErrorContains(t, err, utils.ErrClientNotStarted)

}

// TestManager_SendTransaction_InsufficientFunds tests the SendTransaction method when there are insufficient funds.
// go test -v -cover ./pkg/evm -run TestManager_SendTransaction_InsufficientFunds
func TestManager_SendTransaction_InsufficientFunds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_evm.NewMockClientInterface(ctrl)
	mockClientFactory := mock_evm.NewMockClientFactory(ctrl)
	mockSigner := mock_signer.NewMockTransactionSigner(ctrl)

	// Set up the manager and start the client
	manager := evm.NewManager("http://localhost:8545", mockSigner, mockClientFactory)
	mockClientFactory.EXPECT().DialContext(gomock.Any(), "http://localhost:8545").Return(mockClient, nil)
	manager.Start(context.Background())

	// Create a mock signed transaction
	signedTx := types.NewTransaction(1, common.HexToAddress("0x0"), big.NewInt(1000), 21000, big.NewInt(50), nil)
	signedTxPayload := map[string]interface{}{
		"signedTransaction": signedTx,
	}

	tx := &evm.BaseTransaction{
		TxPayload: signedTxPayload,
	}
	mockSigner.EXPECT().SignTransaction(gomock.Any()).Return(tx, nil)

	clientError := errors.New("insufficient funds for gas * price + value")
	mockClient.EXPECT().SendTransaction(gomock.Any(), signedTx).Return(clientError)

	// Act
	txHash, err := manager.SendTransaction(context.Background(), tx)

	// Assert
	assert.Empty(t, txHash)
	assert.ErrorContains(t, err, utils.ErrEVMInsufficientFunds)
}

// TestManager_SendTransaction_ExceedsBlockGasLimit tests the SendTransaction method when the gas limit exceeds the block's maximum.
// go test -v -cover ./pkg/evm -run TestManager_SendTransaction_ExceedsBlockGasLimit
func TestManager_SendTransaction_ExceedsBlockGasLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_evm.NewMockClientInterface(ctrl)
	mockClientFactory := mock_evm.NewMockClientFactory(ctrl)
	mockSigner := mock_signer.NewMockTransactionSigner(ctrl)

	// Set up the manager and start the client
	manager := evm.NewManager("http://localhost:8545", mockSigner, mockClientFactory)
	mockClientFactory.EXPECT().DialContext(gomock.Any(), "http://localhost:8545").Return(mockClient, nil)
	manager.Start(context.Background())

	// Create a mock signed transaction
	signedTx := types.NewTransaction(1, common.HexToAddress("0x0"), big.NewInt(1000), 21000, big.NewInt(50), nil)
	signedTxPayload := map[string]interface{}{
		"signedTransaction": signedTx,
	}

	tx := &evm.BaseTransaction{
		TxPayload: signedTxPayload,
	}
	mockSigner.EXPECT().SignTransaction(gomock.Any()).Return(tx, nil)

	clientError := errors.New("exceeds block gas limit")
	mockClient.EXPECT().SendTransaction(gomock.Any(), signedTx).Return(clientError)

	// Act
	txHash, err := manager.SendTransaction(context.Background(), tx)

	// Assert
	assert.Empty(t, txHash)
	assert.ErrorContains(t, err, utils.ErrEVMMaxGasCapExceeded)
}

// TestManager_SendTransaction_ReplacementUnderpriced tests the SendTransaction method when the replacement transaction is underpriced.
// go test -v -cover ./pkg/evm -run TestManager_SendTransaction_ReplacementUnderpriced
func TestManager_SendTransaction_ReplacementUnderpriced(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_evm.NewMockClientInterface(ctrl)
	mockClientFactory := mock_evm.NewMockClientFactory(ctrl)
	mockSigner := mock_signer.NewMockTransactionSigner(ctrl)

	// Set up the manager and start the client
	manager := evm.NewManager("http://localhost:8545", mockSigner, mockClientFactory)
	mockClientFactory.EXPECT().DialContext(gomock.Any(), "http://localhost:8545").Return(mockClient, nil)
	manager.Start(context.Background())

	// Create a mock signed transaction
	signedTx := types.NewTransaction(1, common.HexToAddress("0x0"), big.NewInt(1000), 21000, big.NewInt(50), nil)
	signedTxPayload := map[string]interface{}{
		"signedTransaction": signedTx,
	}

	tx := &evm.BaseTransaction{
		TxPayload: signedTxPayload,
	}
	mockSigner.EXPECT().SignTransaction(gomock.Any()).Return(tx, nil)

	clientError := errors.New("replacement transaction underpriced")
	mockClient.EXPECT().SendTransaction(gomock.Any(), signedTx).Return(clientError)

	// Act
	txHash, err := manager.SendTransaction(context.Background(), tx)

	// Assert
	assert.Empty(t, txHash)
	assert.ErrorContains(t, err, utils.ErrEVMReplacementUnderpriced)
}

// TestManager_SendTransaction_NonceTooLow tests the SendTransaction method when the transaction nonce is too low.
// go test -v -cover ./pkg/evm -run TestManager_SendTransaction_NonceTooLow
func TestManager_SendTransaction_NonceTooLow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_evm.NewMockClientInterface(ctrl)
	mockClientFactory := mock_evm.NewMockClientFactory(ctrl)
	mockSigner := mock_signer.NewMockTransactionSigner(ctrl)

	// Set up the manager and start the client
	manager := evm.NewManager("http://localhost:8545", mockSigner, mockClientFactory)
	mockClientFactory.EXPECT().DialContext(gomock.Any(), "http://localhost:8545").Return(mockClient, nil)
	manager.Start(context.Background())

	// Create a mock signed transaction
	signedTx := types.NewTransaction(1, common.HexToAddress("0x0"), big.NewInt(1000), 21000, big.NewInt(50), nil)
	signedTxPayload := map[string]interface{}{
		"signedTransaction": signedTx,
	}

	tx := &evm.BaseTransaction{
		TxPayload: signedTxPayload,
	}
	mockSigner.EXPECT().SignTransaction(gomock.Any()).Return(tx, nil)

	clientError := errors.New("nonce too low")
	mockClient.EXPECT().SendTransaction(gomock.Any(), signedTx).Return(clientError)

	// Act
	txHash, err := manager.SendTransaction(context.Background(), tx)

	// Assert
	assert.Empty(t, txHash)
	assert.ErrorContains(t, err, utils.ErrEVMNonceTooLow)
}
