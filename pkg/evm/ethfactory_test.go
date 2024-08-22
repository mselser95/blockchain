package evm_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mselser95/blockchain/pkg/evm"
	"github.com/stretchr/testify/assert"
)

func TestEthClientFactory_DialContext(t *testing.T) {
	// Arrange
	factory := &evm.EthClientFactory{}
	ctx := context.Background()
	validURL := "https://arb1.arbitrum.io/rpc"
	invalidURL := "invalid_url"

	// Act & Assert

	// Test with a valid URL
	client, err := factory.DialContext(ctx, validURL)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Ensure that the client is of the correct type
	_, ok := client.(*ethclient.Client)
	assert.True(t, ok, "client should be of type *ethclient.Client")

	// Test with an invalid URL
	client, err = factory.DialContext(ctx, invalidURL)
	assert.Error(t, err)
	assert.Nil(t, client)
}
