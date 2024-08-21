package evm_test

import (
	"github.com/mselser95/blockchain/pkg/evm"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mselser95/blockchain/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// To run all tests in this file from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm

// To run this specific test from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm -run TestNewAddress_Success
func TestNewAddress_Success(t *testing.T) {
	validAddress := "0x32Be343B94f860124dC4fEe278FDCBD38C102D88"
	network := utils.Ethereum

	addr, err := evm.NewAddress(validAddress, network)
	assert.NoError(t, err)
	assert.NotNil(t, addr)
	assert.Equal(t, validAddress, addr.String())
	assert.Equal(t, string(network), addr.Network())
}

// To run this specific test from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm -run TestNewAddress_Invalid
func TestNewAddress_Invalid(t *testing.T) {
	invalidAddress := "invalid_address"
	network := utils.Ethereum

	addr, err := evm.NewAddress(invalidAddress, network)
	assert.Error(t, err)
	assert.Nil(t, addr)
}

// To run this specific test from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm -run TestAddress_Validate_Invalid
func TestAddress_Validate_Invalid(t *testing.T) {
	invalidAddress := "0x123"
	network := utils.Ethereum

	addr, err := evm.NewAddress(invalidAddress, network)
	assert.Error(t, err)
	assert.Nil(t, addr)

}

// To run this specific test from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm -run TestAddress_Bytes
func TestAddress_Bytes(t *testing.T) {
	validAddress := "0x32Be343B94f860124dC4fEe278FDCBD38C102D88"
	network := utils.Ethereum

	addr, err := evm.NewAddress(validAddress, network)
	assert.NoError(t, err)
	assert.NotNil(t, addr)

	expectedBytes := common.HexToAddress(validAddress).Bytes()
	assert.Equal(t, expectedBytes, addr.Bytes())
}

// To run this specific test from the root directory with coverage and verbosity:
// go test -v -cover ./pkg/evm -run TestAddress_Network
func TestAddress_Network(t *testing.T) {
	validAddress := "0x32Be343B94f860124dC4fEe278FDCBD38C102D88"
	network := utils.Ethereum

	addr, err := evm.NewAddress(validAddress, network)
	assert.NoError(t, err)
	assert.NotNil(t, addr)

	assert.Equal(t, string(network), addr.Network())
}
