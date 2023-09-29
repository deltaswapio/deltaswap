// SPDX-License-Identifier: Apache 2

pragma solidity ^0.8.19;

import {IDeltaswapRelayer} from "../../interfaces/relayer/IDeltaswapRelayerTyped.sol";

import {getDefaultDeliveryProviderState} from "./DeltaswapRelayerStorage.sol";
import {DeltaswapRelayerGovernance} from "./DeltaswapRelayerGovernance.sol";
import {DeltaswapRelayerSend} from "./DeltaswapRelayerSend.sol";
import {DeltaswapRelayerDelivery} from "./DeltaswapRelayerDelivery.sol";
import {DeltaswapRelayerBase} from "./DeltaswapRelayerBase.sol";

//DeltaswapRelayerGovernance inherits from ERC1967Upgrade, i.e. this is a proxy contract!
contract DeltaswapRelayer is
    DeltaswapRelayerGovernance,
    DeltaswapRelayerSend,
    DeltaswapRelayerDelivery,
    IDeltaswapRelayer
{
    //the only normal storage variable - everything else uses slot pattern
    //no point doing it for this one since it is entirely one-off and of no interest to the rest
    //  of the contract and it also can't accidentally be moved because we are at the bottom of
    //  the inheritance hierarchy here
    bool private initialized;

    constructor(address wormhole) DeltaswapRelayerBase(wormhole) {}

    //needs to be called upon construction of the EC1967 proxy
    function initialize(address defaultDeliveryProvider) public {
        assert(!initialized);
        initialized = true;
        getDefaultDeliveryProviderState().defaultDeliveryProvider = defaultDeliveryProvider;
    }
}
