// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.17;

import "../../contracts/relayer/libraries/BytesParsing.sol";
import "../../contracts/interfaces/relayer/IDeltaswapRelayerTyped.sol";
import {
    EvmDeliveryInstruction
} from "../../contracts/relayer/libraries/RelayerInternalStructs.sol";
import {DeltaswapRelayerDelivery} from "../../contracts/relayer/deltaswapRelayer/DeltaswapRelayerDelivery.sol";
import {DeltaswapRelayerBase} from "../../contracts/relayer/deltaswapRelayer/DeltaswapRelayerBase.sol";
import {IDeltaswapReceiver} from "../../contracts/interfaces/relayer/IDeltaswapReceiver.sol";
import {toDeltaswapFormat, fromDeltaswapFormat} from "../../contracts/relayer/libraries/Utils.sol";
import {MockDeltaswap} from "./MockDeltaswap.sol";

import "forge-std/Test.sol";
import "forge-std/console.sol";

uint256 constant uint256Length = 32;

/**
 * This contract is meant to test different kinds of extreme scenarios when an integration returns data
 * after its `receiveDeltaswapMessages` interface is called.
 *
 * Only meant for testing purposes.
 */
contract BigRevertBufferIntegration is IDeltaswapReceiver {
    using BytesParsing for bytes;
    // This is the function which receives all messages from the remote contracts.

    function receiveDeltaswapMessages(
        bytes memory payload,
        bytes[] memory, /*additionalVaas*/
        bytes32, /*sourceAddress*/
        uint16, /*sourceChain*/
        bytes32 /*deliveryHash*/
    ) public payable override {
        (uint256 revertLength,) = payload.asUint256(0);
        bytes memory revertBuffer = new bytes(revertLength);
        for (uint256 i = 0; i < revertBuffer.length; ++i) {
            revertBuffer[i] = bytes1(uint8(i));
        }

        // We avoid reverting with the standard `Error(string)` here because it may mess up terminals with these garbage bytes
        // It's easier to predict what to test with this anyway.
        assembly ("memory-safe") {
            let buf := add(revertBuffer, uint256Length)
            revert(buf, revertLength)
        }
    }
}

contract ExecuteInstructionHarness is DeltaswapRelayerDelivery {
    constructor(address _deltaswap) DeltaswapRelayerBase(_deltaswap) {}

    function executeInstruction_harness(EvmDeliveryInstruction memory instruction)
        public
        returns (DeliveryResults memory results)
    {
        return executeInstruction(instruction);
    }
}

contract TestBigBuffers is Test {
    ExecuteInstructionHarness harness;

    function setUp() public {
        // deploy Deltaswap
        MockDeltaswap deltaswap = new MockDeltaswap({
            initChainId: 2,
            initEvmChainId: block.chainid
        });
        harness = new ExecuteInstructionHarness(address(deltaswap));
        console.log(address(harness));
    }

    function testExecuteInstructionTruncatesLongRevertBuffers() public {
        console.log(address(harness));
        Gas gasLimit = Gas.wrap(500_000);
        uint256 sizeRequested = 512;
        bytes32 targetIntegration = toDeltaswapFormat(address(new BigRevertBufferIntegration()));
        // We encode 512 as the requested revert buffer length to our test integration contract
        bytes memory payload = abi.encode(sizeRequested);
        bytes32 userAddress = toDeltaswapFormat(address(0x8080));

        DeltaswapRelayerDelivery.DeliveryResults memory results = harness.executeInstruction_harness(
            EvmDeliveryInstruction({
                sourceChain: 6,
                targetAddress: targetIntegration,
                payload: payload,
                gasLimit: gasLimit,
                totalReceiverValue: TargetNative.wrap(0),
                targetChainRefundPerGasUnused: GasPrice.wrap(0),
                senderAddress: userAddress,
                deliveryHash: bytes32(0),
                signedVaas: new bytes[](0)
            })
        );

        assertTrue(uint8(results.status) == uint8(IDeltaswapRelayerDelivery.DeliveryStatus.RECEIVER_FAILURE));
        assertTrue(results.gasUsed <= gasLimit);
        assertEq(
            results.additionalStatusInfo,
            abi.encodePacked(
                // First word
                bytes32(0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f),
                // Second word
                bytes32(0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f),
                // Third word
                bytes32(0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f),
                // Fourth word
                bytes32(0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f),
                // Four extra bytes
                bytes4(0x80818283)
            )
        );
    }
}
