// SPDX-License-Identifier: Apache 2

pragma solidity ^0.8.0;

import "../../contracts/interfaces/relayer/IDeltaswapRelayerTyped.sol";
import {IDeltaswap} from "../../contracts/interfaces/IDeltaswap.sol";
import {DeltaswapSimulator} from "./DeltaswapSimulator.sol";
import {toDeltaswapFormat} from "../../contracts/relayer/libraries/Utils.sol";
import {
    DeliveryInstruction,
    DeliveryOverride,
    RedeliveryInstruction
} from "../../contracts/relayer/libraries/RelayerInternalStructs.sol";
import {DeltaswapRelayerSerde} from
    "../../contracts/relayer/deltaswapRelayer/DeltaswapRelayerSerde.sol";
import "../../contracts/libraries/external/BytesLib.sol";
import "forge-std/Vm.sol";
import "../../contracts/interfaces/relayer/TypedUnits.sol";
import "../../contracts/relayer/libraries/ExecutionParameters.sol";

contract MockGenericRelayer {
    using BytesLib for bytes;
    using WeiLib for Wei;
    using GasLib for Gas;
    using TargetNativeLib for TargetNative;

    IDeltaswap relayerDeltaswap;
    DeltaswapSimulator relayerDeltaswapSimulator;
    uint256 transactionIndex;

    address private constant VM_ADDRESS =
        address(bytes20(uint160(uint256(keccak256("hevm cheat code")))));

    Vm public constant vm = Vm(VM_ADDRESS);

    mapping(uint16 => address) deltaswapRelayerContracts;

    mapping(uint16 => address) relayers;

    mapping(bytes32 => bytes[]) pastEncodedVMs;

    mapping(bytes32 => bytes) pastEncodedDeliveryVAA;

    constructor(address _deltaswap, address _deltaswapSimulator) {
        // deploy Deltaswap

        relayerDeltaswap = IDeltaswap(_deltaswap);
        relayerDeltaswapSimulator = DeltaswapSimulator(_deltaswapSimulator);
        transactionIndex = 0;
    }

    function getPastEncodedVMs(
        uint16 chainId,
        uint64 deliveryVAASequence
    ) public view returns (bytes[] memory) {
        return pastEncodedVMs[keccak256(abi.encodePacked(chainId, deliveryVAASequence))];
    }

    function getPastDeliveryVAA(
        uint16 chainId,
        uint64 deliveryVAASequence
    ) public view returns (bytes memory) {
        return pastEncodedDeliveryVAA[keccak256(abi.encodePacked(chainId, deliveryVAASequence))];
    }

    function setInfo(
        uint16 chainId,
        uint64 deliveryVAASequence,
        bytes[] memory encodedVMs,
        bytes memory encodedDeliveryVAA
    ) internal {
        pastEncodedVMs[keccak256(abi.encodePacked(chainId, deliveryVAASequence))] = encodedVMs;
        pastEncodedDeliveryVAA[keccak256(abi.encodePacked(chainId, deliveryVAASequence))] =
            encodedDeliveryVAA;
    }

    function setDeltaswapRelayerContract(uint16 chainId, address contractAddress) public {
        deltaswapRelayerContracts[chainId] = contractAddress;
    }

    function setProviderDeliveryAddress(uint16 chainId, address deliveryAddress) public {
        relayers[chainId] = deliveryAddress;
    }

    function relay(uint16 chainId) public {
        relay(vm.getRecordedLogs(), chainId, bytes(""));
    }

    function vaaKeyMatchesVAA(
        MessageKey memory messageKey,
        bytes memory signedVaa
    ) internal view returns (bool) {
        if (messageKey.keyType != VAA_KEY_TYPE) {
            return true;
        }
        (VaaKey memory vaaKey,) = DeltaswapRelayerSerde.decodeVaaKey(messageKey.encodedKey, 0);
        return vaaKeyMatchesVAA(vaaKey, signedVaa);
    }

    function vaaKeyMatchesVAA(
        VaaKey memory vaaKey,
        bytes memory signedVaa
    ) internal view returns (bool) {
        IDeltaswap.VM memory parsedVaa = relayerDeltaswap.parseVM(signedVaa);
        return (vaaKey.chainId == parsedVaa.emitterChainId)
            && (vaaKey.emitterAddress == parsedVaa.emitterAddress)
            && (vaaKey.sequence == parsedVaa.sequence);
    }

    function relay(Vm.Log[] memory logs, uint16 chainId, bytes memory deliveryOverrides) public {
        Vm.Log[] memory entries = relayerDeltaswapSimulator.fetchDeltaswapMessageFromLog(logs);
        bytes[] memory encodedVMs = new bytes[](entries.length);
        for (uint256 i = 0; i < encodedVMs.length; i++) {
            encodedVMs[i] = relayerDeltaswapSimulator.fetchSignedMessageFromLogs(
                entries[i], chainId, address(uint160(uint256(bytes32(entries[i].topics[1]))))
            );
        }
        IDeltaswap.VM[] memory parsed = new IDeltaswap.VM[](encodedVMs.length);
        for (uint16 i = 0; i < encodedVMs.length; i++) {
            parsed[i] = relayerDeltaswap.parseVM(encodedVMs[i]);
        }
        for (uint16 i = 0; i < encodedVMs.length; i++) {
            if (
                parsed[i].emitterAddress == toDeltaswapFormat(deltaswapRelayerContracts[chainId])
                    && (parsed[i].emitterChainId == chainId)
            ) {
                genericRelay(encodedVMs[i], encodedVMs, parsed[i], deliveryOverrides);
            }
        }
    }

    function relay(uint16 chainId, bytes memory deliveryOverrides) public {
        relay(vm.getRecordedLogs(), chainId, deliveryOverrides);
    }

    function genericRelay(
        bytes memory encodedDeliveryVAA,
        bytes[] memory encodedVMs,
        IDeltaswap.VM memory parsedDeliveryVAA,
        bytes memory deliveryOverrides
    ) internal {
        uint8 payloadId = parsedDeliveryVAA.payload.toUint8(0);
        if (payloadId == 1) {
            DeliveryInstruction memory instruction =
                DeltaswapRelayerSerde.decodeDeliveryInstruction(parsedDeliveryVAA.payload);

            bytes[] memory encodedVMsToBeDelivered = new bytes[](instruction.messageKeys.length);

            for (uint8 i = 0; i < instruction.messageKeys.length; i++) {
                for (uint8 j = 0; j < encodedVMs.length; j++) {
                    if (vaaKeyMatchesVAA(instruction.messageKeys[i], encodedVMs[j])) {
                        encodedVMsToBeDelivered[i] = encodedVMs[j];
                        break;
                    }
                }
            }

            EvmExecutionInfoV1 memory executionInfo =
                decodeEvmExecutionInfoV1(instruction.encodedExecutionInfo);
            Wei budget = executionInfo.gasLimit.toWei(executionInfo.targetChainRefundPerGasUnused)
                + instruction.requestedReceiverValue.asNative()
                + instruction.extraReceiverValue.asNative();

            uint16 targetChain = instruction.targetChain;

            vm.prank(relayers[targetChain]);
            IDeltaswapRelayerDelivery(deltaswapRelayerContracts[targetChain]).deliver{
                value: budget.unwrap()
            }(
                encodedVMsToBeDelivered,
                encodedDeliveryVAA,
                payable(relayers[targetChain]),
                deliveryOverrides
            );

            setInfo(
                parsedDeliveryVAA.emitterChainId,
                parsedDeliveryVAA.sequence,
                encodedVMsToBeDelivered,
                encodedDeliveryVAA
            );
        } else if (payloadId == 2) {
            RedeliveryInstruction memory instruction =
                DeltaswapRelayerSerde.decodeRedeliveryInstruction(parsedDeliveryVAA.payload);

            DeliveryOverride memory deliveryOverride = DeliveryOverride({
                newExecutionInfo: instruction.newEncodedExecutionInfo,
                newReceiverValue: instruction.newRequestedReceiverValue,
                redeliveryHash: parsedDeliveryVAA.hash
            });

            EvmExecutionInfoV1 memory executionInfo =
                decodeEvmExecutionInfoV1(instruction.newEncodedExecutionInfo);
            Wei budget = executionInfo.gasLimit.toWei(executionInfo.targetChainRefundPerGasUnused)
                + instruction.newRequestedReceiverValue.asNative();

            bytes memory oldEncodedDeliveryVAA = getPastDeliveryVAA(
                instruction.deliveryVaaKey.chainId, instruction.deliveryVaaKey.sequence
            );
            bytes[] memory oldEncodedVMs = getPastEncodedVMs(
                instruction.deliveryVaaKey.chainId, instruction.deliveryVaaKey.sequence
            );

            uint16 targetChain = DeltaswapRelayerSerde.decodeDeliveryInstruction(
                relayerDeltaswap.parseVM(oldEncodedDeliveryVAA).payload
            ).targetChain;

            vm.prank(relayers[targetChain]);
            IDeltaswapRelayerDelivery(deltaswapRelayerContracts[targetChain]).deliver{
                value: budget.unwrap()
            }(
                oldEncodedVMs,
                oldEncodedDeliveryVAA,
                payable(relayers[targetChain]),
                DeltaswapRelayerSerde.encode(deliveryOverride)
            );
        }
    }
}
