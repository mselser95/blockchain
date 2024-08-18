package evm_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mock_evm "github.com/mselser95/blockchain/internal/mock/evm"
	mock_signer "github.com/mselser95/blockchain/internal/mock/signer"
	"github.com/mselser95/blockchain/pkg/evm"
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
	assert.ErrorIs(t, err, utils.ErrAlreadyStarted)
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
	assert.Equal(t, utils.ErrClientNotStarted, err)
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
