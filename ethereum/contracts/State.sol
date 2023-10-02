// contracts/State.sol
// SPDX-License-Identifier: Apache 2

pragma solidity ^0.8.0;

import "./Structs.sol";

contract Events {
    event LogPhylaxSetChanged(
        uint32 oldPhylaxIndex,
        uint32 newPhylaxIndex
    );

    event LogMessagePublished(
        address emitter_address,
        uint32 nonce,
        bytes payload
    );
}

contract Storage {
    struct DeltaswapState {
        Structs.Provider provider;

        // Mapping of phylax_set_index => phylax set
        mapping(uint32 => Structs.PhylaxSet) phylaxSets;

        // Current active phylax set index
        uint32 phylaxSetIndex;

        // Period for which a phylax set stays active after it has been replaced
        uint32 phylaxSetExpiry;

        // Sequence numbers per emitter
        mapping(address => uint64) sequences;

        // Mapping of consumed governance actions
        mapping(bytes32 => bool) consumedGovernanceActions;

        // Mapping of initialized implementations
        mapping(address => bool) initializedImplementations;

        uint256 messageFee;

        // EIP-155 Chain ID
        uint256 evmChainId;
    }
}

contract State {
    Storage.DeltaswapState _state;
}
