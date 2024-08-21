package utils

// Blockchain is an enumeration of supported blockchain networks.
type Blockchain string

const (
	// Ethereum represents the Ethereum blockchain.
	Ethereum Blockchain = "ethereum"
	// Bitcoin represents the Bitcoin blockchain.
	Bitcoin Blockchain = "bitcoin"
	// Solana represents the Solana blockchain.
	Solana Blockchain = "solana"
	// Arbitrum represents the Arbitrum blockchain.
	Arbitrum Blockchain = "arbitrum"
	// Optimism represents the Optimism blockchain.
	Optimism Blockchain = "optimism"
	// Avalanche represents the Avalanche blockchain.
	Avalanche Blockchain = "avalanche"
	// Polygon represents the Polygon blockchain.
	Polygon Blockchain = "polygon"
	// Cosmos represents the Cosmos blockchain.
	Cosmos Blockchain = "cosmos"
	// Bsc represents the Binance Smart Chain blockchain.
	Bsc Blockchain = "bsc"
	// Base represents the Base blockchain.
	Base Blockchain = "base"
	// Blast represents the Blast blockchain.
	Blast Blockchain = "blast"
)
