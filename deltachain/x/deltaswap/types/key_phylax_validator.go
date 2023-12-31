package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// PhylaxValidatorKeyPrefix is the prefix to retrieve all PhylaxValidator
	PhylaxValidatorKeyPrefix = "PhylaxValidator/value/"
)

// PhylaxValidatorKey returns the store key to retrieve a PhylaxValidator from the index fields
func PhylaxValidatorKey(
	phylaxKey []byte,
) []byte {
	var key []byte

	key = append(key, phylaxKey...)
	key = append(key, []byte("/")...)

	return key
}
