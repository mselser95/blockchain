package evm_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mselser95/blockchain/pkg/evm"
	"github.com/mselser95/blockchain/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func generateRandomAddress() utils.Address {
	key, err := crypto.GenerateKey()
	if err != nil {
		panic(err) // panic in test setup if key generation fails
	}
	address := crypto.PubkeyToAddress(key.PublicKey).Hex()
	addr, err := evm.NewAddress(address, utils.Ethereum)
	if err != nil {
		panic(err) // panic in test setup if address creation fails
	}
	return addr
}

func TestBaseTransaction_Validate(t *testing.T) {
	// Arrange
	gasPrice := big.NewInt(50)
	gasLimit := uint64(21000)
	nonce := uint64(1)
	chainId := big.NewInt(1)
	data := []byte{0x0}
	fromAddress := generateRandomAddress()
	toAddress := generateRandomAddress()
	amount := big.NewInt(100)
	txType := utils.Transfer
	status := utils.Pending
	timestamp := time.Now()

	tx := evm.NewTransaction(
		nil, fromAddress, toAddress, amount,
		&txType, &status, &timestamp, nil,
		gasLimit, gasPrice, chainId, nonce, data,
	)

	// Act
	err := tx.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestBaseTransaction_Validate_MissingFrom(t *testing.T) {
	// Arrange
	gasPrice := big.NewInt(50)
	gasLimit := uint64(21000)
	nonce := uint64(1)
	chainId := big.NewInt(1)
	data := []byte{0x0}
	toAddress := generateRandomAddress()
	amount := big.NewInt(100)
	txType := utils.Transfer
	status := utils.Pending
	timestamp := time.Now()

	tx := evm.NewTransaction(
		nil, nil, toAddress, amount,
		&txType, &status, &timestamp, nil,
		gasLimit, gasPrice, chainId, nonce, data,
	)

	// Act
	err := tx.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "missing sender address", err.Error())
}

func TestBaseTransaction_Validate_MissingTo(t *testing.T) {
	// Arrange
	gasPrice := big.NewInt(50)
	gasLimit := uint64(21000)
	nonce := uint64(1)
	chainId := big.NewInt(1)
	data := []byte{0x0}
	fromAddress := generateRandomAddress()
	amount := big.NewInt(100)
	txType := utils.Transfer
	status := utils.Pending
	timestamp := time.Now()

	tx := evm.NewTransaction(
		nil, fromAddress, nil, amount,
		&txType, &status, &timestamp, nil,
		gasLimit, gasPrice, chainId, nonce, data,
	)

	// Act
	err := tx.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "missing recipient address", err.Error())
}

func TestBaseTransaction_Validate_InvalidAmount(t *testing.T) {
	// Arrange
	gasPrice := big.NewInt(50)
	gasLimit := uint64(21000)
	nonce := uint64(1)
	chainId := big.NewInt(1)
	data := []byte{0x0}
	fromAddress := generateRandomAddress()
	toAddress := generateRandomAddress()
	txType := utils.Transfer
	status := utils.Pending
	timestamp := time.Now()

	tx := evm.NewTransaction(
		nil, fromAddress, toAddress, nil,
		&txType, &status, &timestamp, nil,
		gasLimit, gasPrice, chainId, nonce, data,
	)

	// Act
	err := tx.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "invalid transaction amount", err.Error())
}

func TestBaseTransaction_Validate_MissingGasPrice(t *testing.T) {
	// Arrange
	gasLimit := uint64(21000)
	nonce := uint64(1)
	chainId := big.NewInt(1)
	data := []byte{0x0}
	fromAddress := generateRandomAddress()
	toAddress := generateRandomAddress()
	amount := big.NewInt(100)
	txType := utils.Transfer
	status := utils.Pending
	timestamp := time.Now()

	tx := evm.NewTransaction(
		nil, fromAddress, toAddress, amount,
		&txType, &status, &timestamp, nil,
		gasLimit, nil, chainId, nonce, data,
	)

	// Act
	err := tx.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "missing payload 'gasPrice'", err.Error())
}

func TestBaseTransaction_Validate_MissingGasLimit(t *testing.T) {
	// Arrange
	gasPrice := big.NewInt(50)
	nonce := uint64(1)
	chainId := big.NewInt(1)
	data := []byte{0x0}
	fromAddress := generateRandomAddress()
	toAddress := generateRandomAddress()
	amount := big.NewInt(100)
	txType := utils.Transfer
	status := utils.Pending
	timestamp := time.Now()

	tx := evm.NewTransaction(
		nil, fromAddress, toAddress, amount,
		&txType, &status, &timestamp, nil,
		0, gasPrice, chainId, nonce, data,
	)

	// Act
	err := tx.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "missing payload 'gasLimit'", err.Error())
}

func TestBaseTransaction_Validate_MissingNonce(t *testing.T) {
	// Arrange
	gasPrice := big.NewInt(50)
	gasLimit := uint64(21000)
	chainId := big.NewInt(1)
	data := []byte{0x0}
	fromAddress := generateRandomAddress()
	toAddress := generateRandomAddress()
	amount := big.NewInt(100)
	txType := utils.Transfer
	status := utils.Pending
	timestamp := time.Now()

	tx := evm.NewTransaction(
		nil, fromAddress, toAddress, amount,
		&txType, &status, &timestamp, nil,
		gasLimit, gasPrice, chainId, 0, data,
	)

	// Act
	err := tx.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "missing payload 'nonce'", err.Error())
}

func TestBaseTransaction_Validate_MissingChainId(t *testing.T) {
	// Arrange
	gasPrice := big.NewInt(50)
	gasLimit := uint64(21000)
	nonce := uint64(1)
	data := []byte{0x0}
	fromAddress := generateRandomAddress()
	toAddress := generateRandomAddress()
	amount := big.NewInt(100)
	txType := utils.Transfer
	status := utils.Pending
	timestamp := time.Now()

	tx := evm.NewTransaction(
		nil, fromAddress, toAddress, amount,
		&txType, &status, &timestamp, nil,
		gasLimit, gasPrice, nil, nonce, data,
	)

	// Act
	err := tx.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "missing payload 'chainId'", err.Error())
}

func TestBaseTransaction_Validate_InvalidGasLimit(t *testing.T) {
	// Arrange
	gasPrice := big.NewInt(50)
	nonce := uint64(1)
	chainId := big.NewInt(1)
	data := []byte{0x0}
	fromAddress := generateRandomAddress()
	toAddress := generateRandomAddress()
	amount := big.NewInt(100)
	txType := utils.Transfer
	status := utils.Pending
	timestamp := time.Now()

	tx := evm.NewTransaction(
		nil, fromAddress, toAddress, amount,
		&txType, &status, &timestamp, nil,
		uint64(0), gasPrice, chainId, nonce, data,
	)
	// Manually setting the gasLimit with an incorrect type (int instead of uint64)
	tx.SetPayload("gasLimit", 21000)

	// Act
	err := tx.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "payload 'gasLimit' must be uint64", err.Error())
}

func TestBaseTransaction_Validate_InvalidGasPrice(t *testing.T) {
	// Arrange
	gasLimit := uint64(21000)
	nonce := uint64(1)
	chainId := big.NewInt(1)
	data := []byte{0x0}
	fromAddress := generateRandomAddress()
	toAddress := generateRandomAddress()
	amount := big.NewInt(100)
	txType := utils.Transfer
	status := utils.Pending
	timestamp := time.Now()

	tx := evm.NewTransaction(
		nil, fromAddress, toAddress, amount,
		&txType, &status, &timestamp, nil,
		gasLimit, nil, chainId, nonce, data,
	)
	// Manually setting the gasPrice with an incorrect type (int instead of *big.Int)
	tx.SetPayload("gasPrice", 50)

	// Act
	err := tx.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "payload 'gasPrice' must be *big.Int", err.Error())
}

func TestBaseTransaction_Validate_InvalidNonce(t *testing.T) {
	// Arrange
	gasPrice := big.NewInt(50)
	gasLimit := uint64(21000)
	chainId := big.NewInt(1)
	data := []byte{0x0}
	fromAddress := generateRandomAddress()
	toAddress := generateRandomAddress()
	amount := big.NewInt(100)
	txType := utils.Transfer
	status := utils.Pending
	timestamp := time.Now()

	tx := evm.NewTransaction(
		nil, fromAddress, toAddress, amount,
		&txType, &status, &timestamp, nil,
		gasLimit, gasPrice, chainId, 0, data,
	)
	// Manually setting the nonce with an incorrect type (string instead of uint64)
	tx.SetPayload("nonce", "1")

	// Act
	err := tx.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "payload 'nonce' must be uint64", err.Error())
}

func TestBaseTransaction_Validate_InvalidChainId(t *testing.T) {
	// Arrange
	gasPrice := big.NewInt(50)
	gasLimit := uint64(21000)
	nonce := uint64(1)
	data := []byte{0x0}
	fromAddress := generateRandomAddress()
	toAddress := generateRandomAddress()
	amount := big.NewInt(100)
	txType := utils.Transfer
	status := utils.Pending
	timestamp := time.Now()

	tx := evm.NewTransaction(
		nil, fromAddress, toAddress, amount,
		&txType, &status, &timestamp, nil,
		gasLimit, gasPrice, nil, nonce, data,
	)
	// Manually setting the chainId with an incorrect type (int instead of *big.Int)
	tx.SetPayload("chainId", 1)

	// Act
	err := tx.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "payload 'chainId' must be *big.Int", err.Error())
}

func TestBaseTransaction_Validate_InvalidData(t *testing.T) {
	// Arrange
	gasPrice := big.NewInt(50)
	gasLimit := uint64(21000)
	nonce := uint64(1)
	chainId := big.NewInt(1)
	fromAddress := generateRandomAddress()
	toAddress := generateRandomAddress()
	amount := big.NewInt(100)
	txType := utils.Transfer
	status := utils.Pending
	timestamp := time.Now()

	tx := evm.NewTransaction(
		nil, fromAddress, toAddress, amount,
		&txType, &status, &timestamp, nil,
		gasLimit, gasPrice, chainId, nonce, nil,
	)
	// Manually setting the data with an incorrect type (string instead of []byte)
	tx.SetPayload("data", "invalid_data")

	// Act
	err := tx.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "payload 'data' must be a []byte", err.Error())
}
