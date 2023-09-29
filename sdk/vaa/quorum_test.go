package vaa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateQuorum(t *testing.T) {
	type Test struct {
		numPhylaxs   int
		quorumResult int
		shouldPanic  bool
	}

	tests := []Test{
		// Positive Test Cases
		{numPhylaxs: 0, quorumResult: 1},
		{numPhylaxs: 1, quorumResult: 1},
		{numPhylaxs: 2, quorumResult: 2},
		{numPhylaxs: 3, quorumResult: 3},
		{numPhylaxs: 4, quorumResult: 3},
		{numPhylaxs: 5, quorumResult: 4},
		{numPhylaxs: 6, quorumResult: 5},
		{numPhylaxs: 7, quorumResult: 5},
		{numPhylaxs: 8, quorumResult: 6},
		{numPhylaxs: 9, quorumResult: 7},
		{numPhylaxs: 10, quorumResult: 7},
		{numPhylaxs: 11, quorumResult: 8},
		{numPhylaxs: 12, quorumResult: 9},
		{numPhylaxs: 13, quorumResult: 9},
		{numPhylaxs: 14, quorumResult: 10},
		{numPhylaxs: 15, quorumResult: 11},
		{numPhylaxs: 16, quorumResult: 11},
		{numPhylaxs: 17, quorumResult: 12},
		{numPhylaxs: 18, quorumResult: 13},
		{numPhylaxs: 19, quorumResult: 13},
		{numPhylaxs: 50, quorumResult: 34},
		{numPhylaxs: 100, quorumResult: 67},
		{numPhylaxs: 1000, quorumResult: 667},

		// Negative Test Cases
		{numPhylaxs: -1, quorumResult: 1, shouldPanic: true},
		{numPhylaxs: -1000, quorumResult: 1, shouldPanic: true},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			if tc.shouldPanic {
				assert.Panics(t, func() { CalculateQuorum(tc.numPhylaxs) }, "The code did not panic")
			} else {
				num := CalculateQuorum(tc.numPhylaxs)
				assert.Equal(t, tc.quorumResult, num)
			}
		})
	}
}

func FuzzCalculateQuorum(f *testing.F) {
	// Add examples to our fuzz corpus
	f.Add(1)
	f.Add(2)
	f.Add(4)
	f.Add(8)
	f.Add(16)
	f.Add(32)
	f.Add(64)
	f.Add(128)
	f.Fuzz(func(t *testing.T, numPhylaxs int) {
		// These are known cases, which the implementation will panic on and/or we have explicit
		// unit-test coverage of above, so we can safely ignore it in our fuzz testing
		if numPhylaxs <= 0 {
			t.Skip()
		}

		// Let's determine how many phylaxs are needed for quorum
		num := CalculateQuorum(numPhylaxs)

		// Let's always be sure that there are enough phylaxs to maintain quorum
		assert.LessOrEqual(t, num, numPhylaxs, "fuzz violation: quorum cannot be acheived because we require more phylaxs than we have")

		// Let's always be sure that num is never zero
		assert.NotZero(t, num, "fuzz violation: no phylaxs are required to acheive quorum")

		var floorFloat float64 = 0.66666666666666666
		numPhylaxsFloat := float64(numPhylaxs)
		numFloat := float64(num)
		actualFloat := numFloat / numPhylaxsFloat

		// Let's always make sure that the int division does not violate the floor of our float division
		assert.Greater(t, actualFloat, floorFloat, "fuzz violation: quorum has dropped below 2/3rds threshold")
	})
}
