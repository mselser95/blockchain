package utils

// TxHash represents a generic blockchain transaction hash.
//
//go:generate mockgen -destination=../../internal/mock/utils/mock_utils.go -package=mock_utils -source=transaction_hash.go
type TxHash interface {
	// String returns the string representation of the transaction hash.
	String() string

	// Bytes returns the byte representation of the transaction hash.
	Bytes() []byte

	// Validate checks if the transaction hash is valid.
	Validate() error

	// Network returns the network associated with the transaction hash (e.g., mainnet, testnet).
	Network() string
}
