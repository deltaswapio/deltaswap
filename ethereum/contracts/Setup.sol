// contracts/Implementation.sol
// SPDX-License-Identifier: Apache 2

pragma solidity ^0.8.0;
pragma experimental ABIEncoderV2;

import "./Governance.sol";

import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Upgrade.sol";

contract Setup is Setters, ERC1967Upgrade {
    function setup(
        address implementation,
        address[] memory initialPhylaxs,
        uint16 chainId,
        uint16 governanceChainId,
        bytes32 governanceContract,
        uint256 evmChainId
    ) public {
        require(initialPhylaxs.length > 0, "no guardians specified");

        Structs.PhylaxSet memory initialPhylaxSet = Structs.PhylaxSet({
            keys : initialPhylaxs,
            expirationTime : 0
        });

        storePhylaxSet(initialPhylaxSet, 0);
        // initial guardian set index is 0, which is the default value of the storage slot anyways

        setChainId(chainId);

        setGovernanceChainId(governanceChainId);
        setGovernanceContract(governanceContract);

        setEvmChainId(evmChainId);

        _upgradeTo(implementation);

        // See https://github.com/wormhole-foundation/wormhole/issues/1930 for
        // why we set this here
        setInitialized(implementation);
    }
}
