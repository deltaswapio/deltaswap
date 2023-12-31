package types

import (
	"fmt"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PhylaxSetList:        []PhylaxSet{},
		Config:               nil,
		ReplayProtectionList: []ReplayProtection{},
		SequenceCounterList:  []SequenceCounter{},
		ConsensusPhylaxSetIndex: &ConsensusPhylaxSetIndex{
			Index: 0,
		},
		PhylaxValidatorList: []PhylaxValidator{},
		// this line is used by starport scaffolding # genesis/types/default
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated ID in phylaxSet
	phylaxSetIdMap := make(map[uint32]bool)
	for _, elem := range gs.PhylaxSetList {
		if _, ok := phylaxSetIdMap[elem.Index]; ok {
			return fmt.Errorf("duplicated id for phylaxSet")
		}
		phylaxSetIdMap[elem.Index] = true
	}
	// Check for duplicated index in replayProtection
	replayProtectionIndexMap := make(map[string]struct{})

	for _, elem := range gs.ReplayProtectionList {
		index := string(ReplayProtectionKey(elem.Index))
		if _, ok := replayProtectionIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for replayProtection")
		}
		replayProtectionIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in sequenceCounter
	sequenceCounterIndexMap := make(map[string]struct{})

	for _, elem := range gs.SequenceCounterList {
		index := string(SequenceCounterKey(elem.Index))
		if _, ok := sequenceCounterIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for sequenceCounter")
		}
		sequenceCounterIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in phylaxValidator
	phylaxValidatorIndexMap := make(map[string]struct{})

	for _, elem := range gs.PhylaxValidatorList {
		index := string(PhylaxValidatorKey(elem.PhylaxKey))
		if _, ok := phylaxValidatorIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for phylaxValidator")
		}
		phylaxValidatorIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return nil
}
