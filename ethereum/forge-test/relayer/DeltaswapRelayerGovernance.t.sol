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
import {DeliveryProviderStructs} from
    "../../contracts/relayer/deliveryProvider/DeliveryProviderStructs.sol";
import "../../contracts/interfaces/relayer/IDeltaswapRelayerTyped.sol";
import {DeltaswapRelayer} from "../../contracts/relayer/deltaswapRelayer/DeltaswapRelayer.sol";
import {MockGenericRelayer} from "./MockGenericRelayer.sol";
import {MockDeltaswap} from "./MockDeltaswap.sol";
import {IDeltaswap} from "../../contracts/interfaces/IDeltaswap.sol";
import {DeltaswapSimulator, FakeDeltaswapSimulator} from "./DeltaswapSimulator.sol";
import {IDeltaswapReceiver} from "../../contracts/interfaces/relayer/IDeltaswapReceiver.sol";
import {MockRelayerIntegration} from "../../contracts/mock/relayer/MockRelayerIntegration.sol";
import {TestHelpers} from "./TestHelpers.sol";
import {toDeltaswapFormat} from "../../contracts/relayer/libraries/Utils.sol";
import "../../contracts/libraries/external/BytesLib.sol";

import "forge-std/Test.sol";
import "forge-std/console.sol";
import "forge-std/Vm.sol";

contract Brick {
    function checkAndExecuteUpgradeMigration() external view {}
}

contract DeltaswapRelayerGovernanceTests is Test {
    using BytesLib for bytes;

    TestHelpers helpers;

    bytes32 relayerModule = 0x0000000000000000000000000000000000576f726d686f6c6552656c61796572;
    IDeltaswap deltaswap;
    IDeliveryProvider deliveryProvider;
    DeltaswapSimulator deltaswapSimulator;
    IDeltaswapRelayer deltaswapRelayer;

    function setUp() public {
        helpers = new TestHelpers();
        (deltaswap, deltaswapSimulator) = helpers.setUpDeltaswap(1);
        deliveryProvider = helpers.setUpDeliveryProvider(1);
        deltaswapRelayer = helpers.setUpDeltaswapRelayer(deltaswap, address(deliveryProvider));
    }

    struct GovernanceStack {
        bytes message;
        IDeltaswap.VM preSignedMessage;
        bytes signed;
    }

    function signMessage(bytes memory message) internal returns (bytes memory signed) {
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
        signed = deltaswapSimulator.encodeAndSignMessage(preSignedMessage);
    }

    function fillInGovernanceStack(bytes memory message)
        internal
        returns (GovernanceStack memory stack)
    {
        stack.message = message;
        stack.preSignedMessage = IDeltaswap.VM({
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
        stack.signed = deltaswapSimulator.encodeAndSignMessage(stack.preSignedMessage);
    }

    function testSetDefaultDeliveryProvider() public {
        IDeliveryProvider deliveryProviderB = helpers.setUpDeliveryProvider(1);
        IDeliveryProvider deliveryProviderC = helpers.setUpDeliveryProvider(1);

        bytes memory signed = signMessage(
            abi.encodePacked(
                relayerModule,
                uint8(3),
                uint16(1),
                bytes32(uint256(uint160(address(deliveryProviderB))))
            )
        );

        DeltaswapRelayer(payable(address(deltaswapRelayer))).setDefaultDeliveryProvider(signed);

        assertTrue(deltaswapRelayer.getDefaultDeliveryProvider() == address(deliveryProviderB));

        signed = signMessage(
            abi.encodePacked(
                relayerModule,
                uint8(3),
                uint16(1),
                bytes32(uint256(uint160(address(deliveryProviderC))))
            )
        );

        DeltaswapRelayer(payable(address(deltaswapRelayer))).setDefaultDeliveryProvider(signed);

        assertTrue(deltaswapRelayer.getDefaultDeliveryProvider() == address(deliveryProviderC));
    }

    function testRegisterChain() public {
        IDeltaswapRelayer deltaswapRelayer1 =
            helpers.setUpDeltaswapRelayer(deltaswap, address(deliveryProvider));
        IDeltaswapRelayer deltaswapRelayer2 =
            helpers.setUpDeltaswapRelayer(deltaswap, address(deliveryProvider));
        IDeltaswapRelayer deltaswapRelayer3 =
            helpers.setUpDeltaswapRelayer(deltaswap, address(deliveryProvider));

        helpers.registerDeltaswapRelayerContract(
            DeltaswapRelayer(payable(address(deltaswapRelayer1))),
            deltaswap,
            1,
            2,
            toDeltaswapFormat(address(deltaswapRelayer2))
        );

        helpers.registerDeltaswapRelayerContract(
            DeltaswapRelayer(payable(address(deltaswapRelayer1))),
            deltaswap,
            1,
            3,
            toDeltaswapFormat(address(deltaswapRelayer3))
        );

        assertTrue(
            DeltaswapRelayer(payable(address(deltaswapRelayer1))).getRegisteredDeltaswapRelayerContract(
                2
            ) == toDeltaswapFormat(address(deltaswapRelayer2))
        );

        assertTrue(
            DeltaswapRelayer(payable(address(deltaswapRelayer1))).getRegisteredDeltaswapRelayerContract(
                3
            ) == toDeltaswapFormat(address(deltaswapRelayer3))
        );

        vm.expectRevert(
            abi.encodeWithSignature(
                "ChainAlreadyRegistered(uint16,bytes32)",
                3,
                toDeltaswapFormat(address(deltaswapRelayer3))
            )
        );
        helpers.registerDeltaswapRelayerContract(
            DeltaswapRelayer(payable(address(deltaswapRelayer1))),
            deltaswap,
            1,
            3,
            toDeltaswapFormat(address(deltaswapRelayer2))
        );
    }

    function testUpgradeContractToItself() public {
        address payable myDeltaswapRelayer =
            payable(address(helpers.setUpDeltaswapRelayer(deltaswap, address(deliveryProvider))));

        bytes memory noMigrationFunction = signMessage(
            abi.encodePacked(
                relayerModule,
                uint8(2),
                uint16(1),
                toDeltaswapFormat(address(new DeliveryProviderImplementation()))
            )
        );

        vm.expectRevert();
        DeltaswapRelayer(myDeltaswapRelayer).submitContractUpgrade(noMigrationFunction);

        Brick brick = new Brick();
        bytes memory signed = signMessage(
            abi.encodePacked(relayerModule, uint8(2), uint16(1), toDeltaswapFormat(address(brick)))
        );

        DeltaswapRelayer(myDeltaswapRelayer).submitContractUpgrade(signed);

        vm.expectRevert();
        DeltaswapRelayer(myDeltaswapRelayer).getDefaultDeliveryProvider();
    }
}
