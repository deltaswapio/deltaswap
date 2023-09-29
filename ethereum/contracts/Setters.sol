// contracts/Setters.sol
// SPDX-License-Identifier: Apache 2

pragma solidity ^0.8.0;

import "./State.sol";

contract Setters is State {
    function updatePhylaxSetIndex(uint32 newIndex) internal {
        _state.phylaxSetIndex = newIndex;
    }

    function expirePhylaxSet(uint32 index) internal {
        _state.phylaxSets[index].expirationTime = uint32(block.timestamp) + 86400;
    }

    function storePhylaxSet(Structs.PhylaxSet memory set, uint32 index) internal {
        uint setLength = set.keys.length;
        for (uint i = 0; i < setLength; i++) {
            require(set.keys[i] != address(0), "Invalid key");
        }
        _state.phylaxSets[index] = set;
    }

    function setInitialized(address implementatiom) internal {
        _state.initializedImplementations[implementatiom] = true;
    }

    function setGovernanceActionConsumed(bytes32 hash) internal {
        _state.consumedGovernanceActions[hash] = true;
    }

    function setChainId(uint16 chainId) internal {
        _state.provider.chainId = chainId;
    }

    function setGovernanceChainId(uint16 chainId) internal {
        _state.provider.governanceChainId = chainId;
    }

    function setGovernanceContract(bytes32 governanceContract) internal {
        _state.provider.governanceContract = governanceContract;
    }

    function setMessageFee(uint256 newFee) internal {
        _state.messageFee = newFee;
    }

    function setNextSequence(address emitter, uint64 sequence) internal {
        _state.sequences[emitter] = sequence;
    }

    function setEvmChainId(uint256 evmChainId) internal {
        require(evmChainId == block.chainid, "invalid evmChainId");
        _state.evmChainId = evmChainId;
    }
}