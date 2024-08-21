package utils

// Address represents a generic blockchain address.
//
//go:generate mockgen -destination=../../internal/mock/utils/mock_utils.go -package=mock_utils -source=address.go
type Address interface {
	// String returns the string representation of the address.
	String() string

	// Bytes returns the byte representation of the address.
	Bytes() []byte

	// Network returns the network associated with the address (e.g., mainnet, testnet).
	Network() string
}
