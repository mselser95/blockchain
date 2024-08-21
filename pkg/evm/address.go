package evm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/mselser95/blockchain/pkg/utils"
)

import (
	"errors"
)

// Address represents an EVM address.
type Address struct {
	address common.Address
	network utils.Blockchain
}

// NewAddress creates a new EVM address.
func NewAddress(addr string, network utils.Blockchain) (utils.Address, error) {
	if !common.IsHexAddress(addr) {
		return nil, errors.New("invalid EVM address")
	}
	return &Address{
		address: common.HexToAddress(addr),
		network: network,
	}, nil
}

// String returns the string representation of the EVM address.
func (e *Address) String() string {
	return e.address.Hex()
}

// Bytes returns the byte representation of the EVM address.
func (e *Address) Bytes() []byte {
	return e.address.Bytes()
}

// Network returns the network associated with the EVM address.
func (e *Address) Network() string {
	return string(e.network)
}
