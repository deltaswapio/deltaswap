package evm

import (
	"testing"

	"github.com/deltaswapio/deltaswap/sdk/vaa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestPadAddress(t *testing.T) {
	addr := common.HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed")
	expected_addr := vaa.Address{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5a, 0xae, 0xb6, 0x5, 0x3f, 0x3e, 0x94, 0xc9, 0xb9, 0xa0, 0x9f, 0x33, 0x66, 0x94, 0x35, 0xe7, 0xef, 0x1b, 0xea, 0xed}
	assert.Equal(t, PadAddress(addr), expected_addr)
}
