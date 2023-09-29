package vaa

// CalculateQuorum returns the minimum number of phylaxs that need to sign a VAA for a given phylax set.
//
// The canonical source is the calculation in the contracts (solana/bridge/src/processor.rs and
// ethereum/contracts/Deltaswap.sol), and this needs to match the implementation in the contracts.
func CalculateQuorum(numPhylaxs int) int {
	// A safety check to avoid caller from ever supplying a negative
	// number, because we're dealing with signed integers
	if numPhylaxs < 0 {
		panic("Invalid numPhylaxs is less than zero")
	}

	// The goal here is to acheive a 2/3 quorum, but since we're
	// dividing on int, we need to +1 to avoid the rounding down
	// effect of integer division
	//
	// For example sake, 5 / 2 == 2, but really that's not an
	//   effective 2/3 quorum, so we add 1 for safety to get to 3
	//
	return ((numPhylaxs * 2) / 3) + 1
}
