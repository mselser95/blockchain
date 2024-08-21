package evm

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mselser95/blockchain/pkg/utils"
)

const (
	// HashLength is the expected length of the hash
	HashLength = 32
)

// TxHash represents an EVM transaction hash.
type TxHash struct {
	hash    common.Hash
	network string
}

// NewTxHash creates a new EVM transaction hash.
func NewTxHash(hash string, network string) (utils.TxHash, error) {
	if !isHexHash(hash) {
		return nil, errors.New("invalid EVM transaction hash")
	}
	return &TxHash{
		hash:    common.HexToHash(hash),
		network: network,
	}, nil
}

// String returns the string representation of the EVM transaction hash.
func (e *TxHash) String() string {
	return e.hash.Hex()
}

// Bytes returns the byte representation of the EVM transaction hash.
func (e *TxHash) Bytes() []byte {
	return e.hash.Bytes()
}

// Validate checks if the EVM transaction hash is valid.
func (e *TxHash) Validate() error {
	if !isHexHash(e.hash.Hex()) {
		return errors.New("invalid EVM transaction hash")
	}
	return nil
}

// Network returns the network associated with the EVM transaction hash.
func (e *TxHash) Network() string {
	return e.network
}

// isHexHash verifies whether a string can represent a valid hex-encoded
// Ethereum address or not.
func isHexHash(s string) bool {
	if has0xPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*HashLength && isHex(s)
}

// has0xPrefix validates str begins with '0x' or '0X'.
func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// isHexCharacter returns bool of c being a valid hexadecimal.
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// isHex validates whether each byte is valid hexadecimal string.
func isHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}
