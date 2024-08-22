package evm_test

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mselser95/blockchain/pkg/evm"
	"github.com/mselser95/blockchain/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewPrivateKeySigner_Success(t *testing.T) {
	// Generate a new random private key
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)

	// Convert the private key to a hexadecimal string
	privateKeyHex := fmt.Sprintf("%x", crypto.FromECDSA(privateKey))

	signer, err := evm.NewPrivateKeySigner(privateKeyHex)
	assert.NoError(t, err)
	assert.NotNil(t, signer)
}

func TestNewPrivateKeySigner_InvalidKey(t *testing.T) {
	invalidPrivateKey := "invalid_key"

	signer, err := evm.NewPrivateKeySigner(invalidPrivateKey)
	assert.Error(t, err)
	assert.Nil(t, signer)
	assert.ErrorContains(t, err, utils.ErrEVMInvalidPrivateKey)
}

func TestPrivateKeySigner_SignTransaction_Success(t *testing.T) {
	// Generate a new random private key
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)

	// Convert the private key to a hexadecimal string
	privateKeyHex := fmt.Sprintf("%x", crypto.FromECDSA(privateKey))
	signer, err := evm.NewPrivateKeySigner(privateKeyHex)
	assert.NoError(t, err)

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	from, err := evm.NewAddress(fromAddress, utils.Ethereum)
	assert.NoError(t, err)

	toKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	// Derive the public address to the to key
	toAddress := crypto.PubkeyToAddress(toKey.PublicKey).Hex()
	to, err := evm.NewAddress(toAddress, utils.Ethereum)
	assert.NoError(t, err)

	txType := utils.Transfer
	tx := &evm.BaseTransaction{
		FromAddress: from,
		ToAddress:   to,
		TxAmount:    big.NewInt(1000000000000000000),
		TxType:      &txType,
		TxPayload: map[string]interface{}{
			"nonce":    uint64(1),
			"gasPrice": big.NewInt(20000000000),
			"gasLimit": uint64(21000),
			"chainId":  big.NewInt(1),
			"data":     []byte{},
		},
	}

	signedTx, err := signer.SignTransaction(tx)
	assert.NoError(t, err)
	assert.NotNil(t, signedTx)
	assert.NotNil(t, signedTx.SignedTx())

	// Check if the payload contains the signed transaction
	storedSignedTx, ok := signedTx.Payload()["signedTransaction"].(*types.Transaction)
	assert.True(t, ok)
	assert.NotNil(t, storedSignedTx)
}

func TestPrivateKeySigner_SignTransaction_InvalidTransaction(t *testing.T) {
	// Generate a new random private key
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)

	// Convert the private key to a hexadecimal string
	privateKeyHex := fmt.Sprintf("%x", crypto.FromECDSA(privateKey))
	signer, err := evm.NewPrivateKeySigner(privateKeyHex)
	assert.NoError(t, err)

	invalidTx := &evm.BaseTransaction{}

	_, err = signer.SignTransaction(invalidTx)
	assert.Error(t, err)
	assert.ErrorContains(t, err, utils.ErrEVMInvalidTransaction)
}

func TestPrivateKeySigner_SignTransaction_MissingPayloadFields(t *testing.T) {
	// Generate a new random private key
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)

	// Convert the private key to a hexadecimal string
	privateKeyHex := fmt.Sprintf("%x", crypto.FromECDSA(privateKey))
	signer, err := evm.NewPrivateKeySigner(privateKeyHex)
	assert.NoError(t, err)

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	from, err := evm.NewAddress(fromAddress, utils.Ethereum)
	assert.NoError(t, err)

	toKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	// Derive the public address to the to key
	toAddress := crypto.PubkeyToAddress(toKey.PublicKey).Hex()
	to, err := evm.NewAddress(toAddress, utils.Ethereum)
	assert.NoError(t, err)

	txType := utils.Transfer
	tx := &evm.BaseTransaction{
		FromAddress: from,
		ToAddress:   to,
		TxAmount:    big.NewInt(1000000000000000000),
		TxType:      &txType,
		TxPayload: map[string]interface{}{
			"nonce": uint64(1),
			// Missing gasPrice, gasLimit, chainId, and data
		},
	}

	_, err = signer.SignTransaction(tx)
	assert.Error(t, err)
	assert.ErrorContains(t, err, utils.ErrEVMInvalidTransaction)
}
