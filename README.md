# blockchain

Go library for interacting with multiple blockchain networks. Provides a unified interface for balance queries, gas estimation, transaction signing, and on-chain reads across EVM-compatible chains.

## Supported chains

Ethereum, Arbitrum, Optimism, Avalanche, Polygon, BSC, Base, Blast

## Features

- **BlockchainManager** interface — start/stop, get balances (native + ERC20), estimate gas, send transactions, read contract calls
- **TransactionSigner** — pluggable signer abstraction
- **EVM implementation** — full EVM client with ethclient, transaction building, event log parsing
- Supports native tokens, ERC20, SPL tokens, and Cosmos denoms

## Usage

```go
import "github.com/mselser95/blockchain/pkg/manager"

mgr := evm.NewManager(config)
mgr.Start(ctx)
balance, err := mgr.GetBalance(ctx, address)
```

## Development

```bash
make test       # run tests
make mockgen    # regenerate mocks
```
