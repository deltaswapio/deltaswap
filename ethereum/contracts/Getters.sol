// contracts/Getters.sol
// SPDX-License-Identifier: Apache 2

pragma solidity ^0.8.0;

import "./State.sol";

contract Getters is State {
    function getPhylaxSet(uint32 index) public view returns (Structs.PhylaxSet memory) {
        return _state.phylaxSets[index];
    }

    function getCurrentPhylaxSetIndex() public view returns (uint32) {
        return _state.phylaxSetIndex;
    }

    function getPhylaxSetExpiry() public view returns (uint32) {
        return _state.phylaxSetExpiry;
    }

    function governanceActionIsConsumed(bytes32 hash) public view returns (bool) {
        return _state.consumedGovernanceActions[hash];
    }

    function isInitialized(address impl) public view returns (bool) {
        return _state.initializedImplementations[impl];
    }

    function chainId() public view returns (uint16) {
        return _state.provider.chainId;
    }

    function evmChainId() public view returns (uint256) {
        return _state.evmChainId;
    }

    function isFork() public view returns (bool) {
        return evmChainId() != block.chainid;
    }

    function governanceChainId() public view returns (uint16){
        return _state.provider.governanceChainId;
    }

    function governanceContract() public view returns (bytes32){
        return _state.provider.governanceContract;
    }

    function messageFee() public view returns (uint256) {
        return _state.messageFee;
    }

    function nextSequence(address emitter) public view returns (uint64) {
        return _state.sequences[emitter];
    }
}