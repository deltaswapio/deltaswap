package devnet

import (
	"fmt"

	"github.com/deltaswapio/deltaswap/node/pkg/common"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

// GenerateAndStoreDevnetGuardianKey returns a deterministic testnet key.
func GenerateAndStoreDevnetPhylaxKey(filename string) error {
	// Figure out our devnet index
	idx, err := GetDevnetIndex()
	if err != nil {
		return err
	}

	// Generate the guardian key.
	gk := InsecureDeterministicEcdsaKeyByIndex(ethcrypto.S256(), uint64(idx))

	// Store it to disk.
	if err := common.WriteArmoredKey(gk, "auto-generated deterministic devnet key", filename, common.PhylaxKeyArmoredBlock, true); err != nil {
		return fmt.Errorf("failed to store generated guardian key: %w", err)
	}

	return nil
}
