// SPDX-License-Identifier: Apache 2

pragma solidity ^0.8.0;

import "contracts/interfaces/IDeltaswap.sol";

interface IMyDeltaswap is IDeltaswap {

    function getImplementation() external returns (address);
    function upgradeImpl(address newImplementation) external;
}
