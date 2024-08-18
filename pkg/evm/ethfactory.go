package evm

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
)

// ClientFactory defines an interface for creating a new ClientInterface.
//
//go:generate mockgen -destination=../../internal/mock/evm/mock_factory.go -package=mock_evm -source=ethfactory.go
type ClientFactory interface {
	DialContext(ctx context.Context, url string) (ClientInterface, error)
}

// EthClientFactory is a concrete implementation of ClientFactory that uses ethclient.
type EthClientFactory struct{}

// DialContext dials a new Ethereum client.
func (f *EthClientFactory) DialContext(ctx context.Context, url string) (ClientInterface, error) {
	return ethclient.DialContext(ctx, url)
}
