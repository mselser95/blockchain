package evm

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mselser95/blockchain/pkg/signer"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mselser95/blockchain/pkg/utils"
)

// PrivateKeySigner is a struct that holds a private key and implements the TransactionSigner interface.
type PrivateKeySigner struct {
	privateKey *ecdsa.PrivateKey
}

// NewPrivateKeySigner creates a new PrivateKeySigner instance with the provided private key.
func NewPrivateKeySigner(privateKey string) (signer.TransactionSigner, error) {
	// Decode the private key
	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, utils.WrapError(utils.ErrEVMInvalidPrivateKey, err)
	}

	return &PrivateKeySigner{
		privateKey: pk,
	}, nil
}

// SignTransaction signs the provided transaction with the private key.
func (pks *PrivateKeySigner) SignTransaction(tx utils.Transaction) (utils.Transaction, error) {
	if tx == nil {
		return nil, utils.WrapError(utils.ErrEVMInvalidTransaction)
	}

	if err := tx.Validate(); err != nil {
		return nil, utils.WrapError(utils.ErrEVMInvalidTransaction, err)
	}

	// Prepare the transaction parameters from the utils.Transaction
	nonce := tx.Payload()["nonce"].(uint64)
	gasPrice := tx.Payload()["gasPrice"].(*big.Int)
	gasLimit := tx.Payload()["gasLimit"].(uint64)
	chainId := tx.Payload()["chainId"].(*big.Int)
	to := common.HexToAddress(tx.To().String())
	value := tx.Amount()
	data := tx.Payload()["data"].([]byte)

	// Create the transaction object
	signedTx := types.NewTx(&types.LegacyTx{
		Nonce: nonce, GasPrice: gasPrice, Gas: gasLimit, To: &to, Value: value, Data: data,
	})

	// Sign the transaction using the private key
	txSigner := types.LatestSignerForChainID(chainId)
	signedTransaction, err := types.SignTx(signedTx, txSigner, pks.privateKey)
	if err != nil {
		return nil, utils.WrapError(utils.ErrEVMFailedToSignTransaction, err)
	}

	// Serialize the signed transaction to bytes
	txBytes, err := signedTransaction.MarshalBinary()
	if err != nil {
		return nil, utils.WrapError(utils.ErrEVMFailedToSignTransaction, err)
	}

	// Set the signed transaction bytes in the original transaction
	tx.SetSignedTx(txBytes)
	tx.SetPayload("signedTransaction", signedTransaction)

	return tx, nil
}
