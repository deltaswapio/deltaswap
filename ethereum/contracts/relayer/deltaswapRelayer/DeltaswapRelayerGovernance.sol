// SPDX-License-Identifier: Apache 2
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Upgrade.sol";

import {IDeltaswap} from "../../interfaces/IDeltaswap.sol";
import {InvalidPayloadLength} from "../../interfaces/relayer/IDeltaswapRelayerTyped.sol";
import {fromDeltaswapFormat} from "../../relayer/libraries/Utils.sol";
import {BytesParsing} from "../../relayer/libraries/BytesParsing.sol";
import {
    getGovernanceState,
    getRegisteredDeltaswapRelayersState,
    getDefaultDeliveryProviderState
} from "./DeltaswapRelayerStorage.sol";
import {DeltaswapRelayerBase} from "./DeltaswapRelayerBase.sol";

error GovernanceActionAlreadyConsumed(bytes32 hash);
error InvalidGovernanceVM(string reason);
error InvalidGovernanceChainId(uint16 parsed, uint16 expected);
error InvalidGovernanceContract(bytes32 parsed, bytes32 expected);

error InvalidPayloadChainId(uint16 parsed, uint16 expected);
error InvalidPayloadAction(uint8 parsed, uint8 expected);
error InvalidPayloadModule(bytes32 parsed, bytes32 expected);
error InvalidFork();
error ContractUpgradeFailed(bytes failure);
error ChainAlreadyRegistered(uint16 chainId, bytes32 registeredDeltaswapRelayerContract);
error InvalidDefaultDeliveryProvider(bytes32 defaultDeliveryProvider);

abstract contract DeltaswapRelayerGovernance is DeltaswapRelayerBase, ERC1967Upgrade {
    //This constant should actually be defined in IDeltaswap. Alas, it isn't.
    uint16 private constant DELTASWAP_CHAINID_UNSET = 0;

    /**
     * Governance VMs are encoded in a packed fashion using the general deltaswap scheme:
     *   GovernancePacket = <Common Header|Action Parameters>
     *
     * For a more detailed explanation see here:
     *   - https://docs.deltaswap.com/deltaswap/governance
     *   - https://github.com/deltaswapio/deltaswap/blob/main/whitepapers/0002_governance_messaging.md
     */

    //Right shifted ascii encoding of "WormholeRelayer"
    bytes32 private constant module =
        0x0000000000000000000000000000000000576f726d686f6c6552656c61796572;

    /**
     * The choice of action enumeration and parameters follows the scheme of the core bridge:
     *   - https://github.com/deltaswapio/deltaswap/blob/main/ethereum/contracts/bridge/BridgeGovernance.sol#L115
     */

    /**
     * Registers a deltaswap relayer contract that was deployed on another chain with the DeltaswapRelayer on
     *   this chain. The equivalent to the core bridge's registerChain action.
     *
     * Action Parameters:
     *   - uint16 foreignChainId
     *   - bytes32 foreignContractAddress
     */
    uint8 private constant GOVERNANCE_ACTION_REGISTER_DELTASWAP_RELAYER_CONTRACT = 1;

    /**
     * Upgrades the DeltaswapRelayer contract to a new implementation. The equivalent to the core bridge's
     *   upgrade action.
     *
     * Action Parameters:
     *   - bytes32 newImplementation
     */
    uint8 private constant GOVERNANCE_ACTION_CONTRACT_UPGRADE = 2;

    /**
     * Sets the default relay provider for the DeltaswapRelayer. Has no equivalent in the core bridge.
     *
     * Action Parameters:
     *   - bytes32 newProvider
     */
    uint8 private constant GOVERNANCE_ACTION_UPDATE_DEFAULT_PROVIDER = 3;

    //By checking that only the contract can call itself, we can enforce that the migration code is
    //  executed upon program upgrade and that it can't be called externally by anyone else.
    function checkAndExecuteUpgradeMigration() external {
        assert(msg.sender == address(this));
        executeUpgradeMigration();
    }

    function executeUpgradeMigration() internal virtual {
        //override and implement in DeltaswapRelayer upon contract upgrade (if required)
    }

    function registerDeltaswapRelayerContract(bytes memory encodedVm) external {
        (uint16 foreignChainId, bytes32 foreignAddress) =
            parseAndCheckRegisterDeltaswapRelayerContractVm(encodedVm);

        getRegisteredDeltaswapRelayersState().registeredDeltaswapRelayers[foreignChainId] =
            foreignAddress;
    }

    event ContractUpgraded(address indexed oldContract, address indexed newContract);

    function submitContractUpgrade(bytes memory encodedVm) external {
        address currentImplementation = _getImplementation();
        address newImplementation = parseAndCheckContractUpgradeVm(encodedVm);

        _upgradeTo(newImplementation);

        (bool success, bytes memory revertData) =
            address(this).call(abi.encodeCall(this.checkAndExecuteUpgradeMigration, ()));

        if (!success) {
            revert ContractUpgradeFailed(revertData);
        }

        emit ContractUpgraded(currentImplementation, newImplementation);
    }

    function setDefaultDeliveryProvider(bytes memory encodedVm) external {
        address newProvider = parseAndCheckRegisterDefaultDeliveryProviderVm(encodedVm);

        getDefaultDeliveryProviderState().defaultDeliveryProvider = newProvider;
    }

    // ------------------------------------------- PRIVATE -------------------------------------------
    using BytesParsing for bytes;

    function parseAndCheckRegisterDeltaswapRelayerContractVm(bytes memory encodedVm)
        private
        returns (uint16 foreignChainId, bytes32 foreignAddress)
    {
        bytes memory payload = verifyAndConsumeGovernanceVM(encodedVm);
        uint256 offset = parseAndCheckPayloadHeader(
            payload, GOVERNANCE_ACTION_REGISTER_DELTASWAP_RELAYER_CONTRACT, true
        );

        (foreignChainId, offset) = payload.asUint16Unchecked(offset);
        (foreignAddress, offset) = payload.asBytes32Unchecked(offset);

        checkLength(payload, offset);

        if (getRegisteredDeltaswapRelayerContract(foreignChainId) != bytes32(0)) {
            revert ChainAlreadyRegistered(
                foreignChainId, getRegisteredDeltaswapRelayerContract(foreignChainId)
            );
        }
    }

    function parseAndCheckContractUpgradeVm(bytes memory encodedVm)
        private
        returns (address newImplementation)
    {
        bytes memory payload = verifyAndConsumeGovernanceVM(encodedVm);
        uint256 offset =
            parseAndCheckPayloadHeader(payload, GOVERNANCE_ACTION_CONTRACT_UPGRADE, false);

        bytes32 newImplementationWhFmt;
        (newImplementationWhFmt, offset) = payload.asBytes32Unchecked(offset);
        //fromDeltaswapFormat reverts if first 12 bytes aren't zero (i.e. if it's not an EVM address)
        newImplementation = fromDeltaswapFormat(newImplementationWhFmt);

        checkLength(payload, offset);
    }

    function parseAndCheckRegisterDefaultDeliveryProviderVm(bytes memory encodedVm)
        private
        returns (address newProvider)
    {
        bytes memory payload = verifyAndConsumeGovernanceVM(encodedVm);
        uint256 offset =
            parseAndCheckPayloadHeader(payload, GOVERNANCE_ACTION_UPDATE_DEFAULT_PROVIDER, false);

        bytes32 newProviderWhFmt;
        (newProviderWhFmt, offset) = payload.asBytes32Unchecked(offset);
        //fromDeltaswapFormat reverts if first 12 bytes aren't zero (i.e. if it's not an EVM address)
        newProvider = fromDeltaswapFormat(newProviderWhFmt);

        checkLength(payload, offset);

        if (newProvider == address(0)) {
            revert InvalidDefaultDeliveryProvider(newProviderWhFmt);
        }
    }

    function verifyAndConsumeGovernanceVM(bytes memory encodedVm)
        private
        returns (bytes memory payload)
    {
        (IDeltaswap.VM memory vm, bool valid, string memory reason) =
            getDeltaswap().parseAndVerifyVM(encodedVm);

        if (!valid) {
            revert InvalidGovernanceVM(reason);
        }

        uint16 governanceChainId = getDeltaswap().governanceChainId();
        if (vm.emitterChainId != governanceChainId) {
            revert InvalidGovernanceChainId(vm.emitterChainId, governanceChainId);
        }

        bytes32 governanceContract = getDeltaswap().governanceContract();
        if (vm.emitterAddress != governanceContract) {
            revert InvalidGovernanceContract(vm.emitterAddress, governanceContract);
        }

        bool consumed = getGovernanceState().consumedGovernanceActions[vm.hash];
        if (consumed) {
            revert GovernanceActionAlreadyConsumed(vm.hash);
        }

        getGovernanceState().consumedGovernanceActions[vm.hash] = true;

        return vm.payload;
    }

    function parseAndCheckPayloadHeader(
        bytes memory encodedPayload,
        uint8 expectedAction,
        bool allowUnset
    ) private view returns (uint256 offset) {
        bytes32 parsedModule;
        (parsedModule, offset) = encodedPayload.asBytes32Unchecked(offset);
        if (parsedModule != module) {
            revert InvalidPayloadModule(parsedModule, module);
        }

        uint8 parsedAction;
        (parsedAction, offset) = encodedPayload.asUint8Unchecked(offset);
        if (parsedAction != expectedAction) {
            revert InvalidPayloadAction(parsedAction, expectedAction);
        }

        uint16 parsedChainId;
        (parsedChainId, offset) = encodedPayload.asUint16Unchecked(offset);
        if (!(parsedChainId == DELTASWAP_CHAINID_UNSET && allowUnset)) {
            if (getDeltaswap().isFork()) {
                revert InvalidFork();
            }

            if (parsedChainId != getChainId()) {
                revert InvalidPayloadChainId(parsedChainId, getChainId());
            }
        }
    }

    function checkLength(bytes memory payload, uint256 expected) private pure {
        if (payload.length != expected) {
            revert InvalidPayloadLength(payload.length, expected);
        }
    }
}
