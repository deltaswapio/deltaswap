// SPDX-License-Identifier: Apache 2

pragma solidity ^0.8.0;
import "contracts/Setters.sol";

contract MySetters is Setters {

    function updatePhylaxSetIndex_external(uint32 newIndex) external {
        updatePhylaxSetIndex(newIndex);
    }

    function expirePhylaxSet_external(uint32 index) external {
        expirePhylaxSet(index);
    }

    function setInitialized_external(address implementation) external {
        setInitialized(implementation);
    }

    function setGovernanceActionConsumed_external(bytes32 hash) external {
        setGovernanceActionConsumed(hash);
    }

    function setChainId_external(uint16 chainId) external {
        setChainId(chainId);
    }

    function setGovernanceChainId_external(uint16 chainId) external {
        setGovernanceChainId(chainId);
    }

    function setGovernanceContract_external(bytes32 governanceContract) external {
        setGovernanceContract(governanceContract);
    }

    function setMessageFee_external(uint256 newFee) external {
        setMessageFee(newFee);
    }

    function setNextSequence_external(address emitter, uint64 sequence) external {
        setNextSequence(emitter, sequence);
    }

    function setEvmChainId_external(uint256 evmChainId) external {
        setEvmChainId(evmChainId);
    }
}
