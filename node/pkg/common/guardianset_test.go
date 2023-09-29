package common

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestKeyIndex(t *testing.T) {
	type test struct {
		phylaxSet PhylaxSet
		address   string
		result    bool
		keyIndex  int
	}

	phylaxSet := PhylaxSet{
		Keys: []common.Address{
			common.HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed"),
			common.HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaee"),
		},
		Index: 1,
	}

	tests := []test{
		{phylaxSet: phylaxSet, address: "0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed", result: true, keyIndex: 0},
		{phylaxSet: phylaxSet, address: "0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaee", result: true, keyIndex: 1},
		{phylaxSet: phylaxSet, address: "0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaef", result: false, keyIndex: -1},
	}

	for _, testCase := range tests {
		t.Run(testCase.address, func(t *testing.T) {
			gs := testCase.phylaxSet
			keyIndex, result := gs.KeyIndex(common.HexToAddress(testCase.address))
			assert.Equal(t, result, testCase.result)
			assert.Equal(t, keyIndex, testCase.keyIndex)
		})
	}
}

func TestKeysAsHexStrings(t *testing.T) {
	gs := PhylaxSet{
		Keys: []common.Address{
			common.HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed"),
			common.HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaee"),
		},
		Index: 1,
	}

	keyStrings := make([]string, len(gs.Keys))
	// Note that go-ethereum common.HexToAddress() will force a supplied all-lower
	// address to mixedcase for a valid checksum, which is why these are different
	keyStrings[0] = "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed"
	keyStrings[1] = "0x5Aaeb6053f3e94c9B9a09F33669435e7EF1BeAee"
	assert.Equal(t, keyStrings, gs.KeysAsHexStrings())
}

func TestNewPhylaxSetState(t *testing.T) {
	gss := NewPhylaxSetState(nil)
	assert.NotNil(t, gss)
	assert.Nil(t, gss.current)
	assert.Nil(t, gss.Get())
}

func TestSet(t *testing.T) {
	var gs PhylaxSet = PhylaxSet{
		Keys: []common.Address{
			common.HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed"),
			common.HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaee"),
		},
		Index: 1,
	}

	gss := NewPhylaxSetState(nil)
	assert.Nil(t, gss.current)
	gss.Set(&gs)
	assert.Equal(t, gss.current, &gs)
}

func TestGet(t *testing.T) {
	var gs PhylaxSet = PhylaxSet{
		Keys: []common.Address{
			common.HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed"),
			common.HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaee"),
		},
		Index: 1,
	}

	gss := NewPhylaxSetState(nil)
	assert.Nil(t, gss.Get())
	gss.Set(&gs)
	assert.Equal(t, gss.Get(), &gs)
}
