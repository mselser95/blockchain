package evm_test

import (
	"github.com/mselser95/blockchain/pkg/evm"

	"github.com/stretchr/testify/assert"
	"testing"
)

// To run all tests in this file from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm

// To run this specific test from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm -run TestNewTxHash_Success
func TestNewTxHash_Success(t *testing.T) {
	validHash := "0x5d5f3c8b4692ff7e93556f96d9b11e1a2b22305a109b5d0c0848b72afdf11650"
	network := "mainnet"

	txHash, err := evm.NewTxHash(validHash, network)
	assert.NoError(t, err)
	assert.NotNil(t, txHash)
	assert.Equal(t, validHash, txHash.String())
	assert.Equal(t, network, txHash.Network())
}

// To run this specific test from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm -run TestNewTxHash_InvalidHash
func TestNewTxHash_InvalidHash(t *testing.T) {
	invalidHash := "invalidhash"
	network := "mainnet"

	txHash, err := evm.NewTxHash(invalidHash, network)
	assert.Error(t, err)
	assert.Nil(t, txHash)
}

// To run this specific test from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm -run TestTxHash_Validate_Success
func TestTxHash_Validate_Success(t *testing.T) {
	validHash := "0x5d5f3c8b4692ff7e93556f96d9b11e1a2b22305a109b5d0c0848b72afdf11650"
	network := "mainnet"

	txHash, err := evm.NewTxHash(validHash, network)
	assert.NoError(t, err)
	assert.NotNil(t, txHash)

	err = txHash.Validate()
	assert.NoError(t, err)
}

// To run this specific test from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm -run TestTxHash_Validate_Invalid
func TestTxHash_Validate_Invalid(t *testing.T) {
	invalidHash := "0x123"
	network := "mainnet"

	txHash, err := evm.NewTxHash(invalidHash, network)
	assert.Error(t, err)
	assert.Nil(t, txHash)
}

// To run this specific test from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm -run TestTxHash_Bytes
func TestTxHash_Bytes(t *testing.T) {
	validHash := "0x5d5f3c8b4692ff7e93556f96d9b11e1a2b22305a109b5d0c0848b72afdf11650"
	network := "mainnet"

	txHash, err := evm.NewTxHash(validHash, network)
	assert.NoError(t, err)
	assert.NotNil(t, txHash)

	expectedBytes := []byte{93, 95, 60, 139, 70, 146, 255, 126, 147, 85, 111, 150, 217, 177, 30, 26, 43, 34, 48, 90, 16, 155, 93, 12, 8, 72, 183, 42, 253, 241, 22, 80}
	assert.Equal(t, expectedBytes, txHash.Bytes())
}

// To run this specific test from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm -run TestTxHash_Network
func TestTxHash_Network(t *testing.T) {
	validHash := "0x5d5f3c8b4692ff7e93556f96d9b11e1a2b22305a109b5d0c0848b72afdf11650"
	network := "testnet"

	txHash, err := evm.NewTxHash(validHash, network)
	assert.NoError(t, err)
	assert.NotNil(t, txHash)

	assert.Equal(t, network, txHash.Network())
}
