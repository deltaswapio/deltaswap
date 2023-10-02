// SPDX-License-Identifier: Apache 2

pragma solidity ^0.8.0;

import {IDeliveryProvider} from "../../contracts/interfaces/relayer/IDeliveryProviderTyped.sol";
import {DeliveryProvider} from "../../contracts/relayer/deliveryProvider/DeliveryProvider.sol";
import {DeliveryProviderSetup} from
    "../../contracts/relayer/deliveryProvider/DeliveryProviderSetup.sol";
import {DeliveryProviderImplementation} from
    "../../contracts/relayer/deliveryProvider/DeliveryProviderImplementation.sol";
import {DeliveryProviderProxy} from
    "../../contracts/relayer/deliveryProvider/DeliveryProviderProxy.sol";
import {IDeltaswapRelayer} from "../../contracts/interfaces/relayer/IDeltaswapRelayerTyped.sol";
import {DeltaswapRelayer} from "../../contracts/relayer/deltaswapRelayer/DeltaswapRelayer.sol";
import {Create2Factory} from "../../contracts/relayer/create2Factory/Create2Factory.sol";
import {MockGenericRelayer} from "./MockGenericRelayer.sol";
import {MockDeltaswap} from "./MockDeltaswap.sol";
import {IDeltaswap} from "../../contracts/interfaces/IDeltaswap.sol";
import {DeltaswapSimulator, FakeDeltaswapSimulator} from "./DeltaswapSimulator.sol";
import "../../contracts/libraries/external/BytesLib.sol";

import "forge-std/Test.sol";
import "forge-std/console.sol";
import "forge-std/Vm.sol";

contract TestHelpers {
    using BytesLib for bytes;

    address private constant VM_ADDRESS =
        address(bytes20(uint160(uint256(keccak256("hevm cheat code")))));

    Vm public constant vm = Vm(VM_ADDRESS);

    DeltaswapSimulator helperDeltaswapSimulator;

    constructor() {
        (, helperDeltaswapSimulator) = setUpDeltaswap(1);
    }

    function registerDeltaswapRelayerContract(
        DeltaswapRelayer governance,
        IDeltaswap deltaswap,
        uint16 currentChainId,
        uint16 chainId,
        bytes32 coreRelayerContractAddress
    ) public {
        bytes32 deltaswapRelayerModule =
            0x0000000000000000000000000000000000576f726d686f6c6552656c61796572;
        bytes memory message = abi.encodePacked(
            deltaswapRelayerModule, uint8(1), currentChainId, chainId, coreRelayerContractAddress
        );
        IDeltaswap.VM memory preSignedMessage = IDeltaswap.VM({
            version: 1,
            timestamp: uint32(block.timestamp),
            nonce: 0,
            emitterChainId: deltaswap.governanceChainId(),
            emitterAddress: deltaswap.governanceContract(),
            sequence: 0,
            consistencyLevel: 200,
            payload: message,
            phylaxSetIndex: 0,
            signatures: new IDeltaswap.Signature[](0),
            hash: bytes32("")
        });

        bytes memory signed = helperDeltaswapSimulator.encodeAndSignMessage(preSignedMessage);
        governance.registerDeltaswapRelayerContract(signed);
    }

    function setUpDeltaswap(uint16 chainId)
        public
        returns (IDeltaswap deltaswapContract, DeltaswapSimulator deltaswapSimulator)
    {
        // deploy Deltaswap
        MockDeltaswap deltaswap = new MockDeltaswap({
            initChainId: chainId,
            initEvmChainId: block.chainid
        });

        // replace Deltaswap with the Deltaswap Simulator contract (giving access to some nice helper methods for signing)
        deltaswapSimulator = new FakeDeltaswapSimulator(
            deltaswap
        );

        deltaswapContract = deltaswap;
    }

    function setUpDeliveryProvider(
        uint16 chainId
    ) public returns (DeliveryProvider deliveryProvider) {
        vm.prank(msg.sender);
        DeliveryProviderSetup deliveryProviderSetup = new DeliveryProviderSetup();
        vm.prank(msg.sender);
        DeliveryProviderImplementation deliveryProviderImplementation =
            new DeliveryProviderImplementation();
        vm.prank(msg.sender);
        DeliveryProviderProxy myDeliveryProvider = new DeliveryProviderProxy(
            address(deliveryProviderSetup),
            abi.encodeCall(
                DeliveryProviderSetup.setup,
                (
                    address(deliveryProviderImplementation),
                    chainId
                )
            )
        );

        deliveryProvider = DeliveryProvider(address(myDeliveryProvider));
    }

    function setUpDeltaswapRelayer(
        IDeltaswap deltaswap,
        address defaultDeliveryProvider
    ) public returns (IDeltaswapRelayer coreRelayer) {
        Create2Factory create2Factory = new Create2Factory();

        address proxyAddressComputed =
            create2Factory.computeProxyAddress(address(this), "0xGenericRelayer");

        DeltaswapRelayer coreRelayerImplementation = new DeltaswapRelayer(address(deltaswap));

        bytes memory initCall =
            abi.encodeCall(DeltaswapRelayer.initialize, (defaultDeliveryProvider));

        coreRelayer = IDeltaswapRelayer(
            create2Factory.create2Proxy(
                "0xGenericRelayer", address(coreRelayerImplementation), initCall
            )
        );
        require(
            address(coreRelayer) == proxyAddressComputed, "computed must match actual proxy addr"
        );
    }
}
