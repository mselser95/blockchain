package utils

// TokenType is an enumeration of supported token types.
type TokenType int

const (
	// Native represents a native token on the blockchain network.
	Native TokenType = iota
	// ERC20 represents an ERC20 token on the Ethereum network.
	ERC20
	// SPLToken represents a SPL token on the Solana network.
	SPLToken
	// CosmosDenom represents a Cosmos denomination on the Cosmos network.
	CosmosDenom
)

// Token represents an abstract token on a blockchain network.
type Token struct {
	// Type of the token (e.g., Native, ERC20, SPLToken, CosmosDenom)
	Type TokenType
	// Address or identifier of the token (contract address, mint address, denomination)
	Address *Address
	// Human-readable name of the token
	Name string
	// Symbol of the token (e.g., ETH, SOL, ATOM)
	Symbol string
	// Number of decimals used by the token
	Decimals int
}
