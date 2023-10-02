// SPDX-License-Identifier: Apache 2
pragma solidity ^0.8.0;

import {IDeltaswap} from "../../contracts/interfaces/IDeltaswap.sol";
import {MockDeltaswap} from "./MockDeltaswap.sol";
import "../../contracts/libraries/external/BytesLib.sol";

import "forge-std/Vm.sol";
import "forge-std/console.sol";

/**
 * @notice These are the common parts for the signing and the non signing deltaswap simulators.
 * @dev This contract is meant to be used when testing against a mainnet fork.
 */
abstract contract DeltaswapSimulator {
    using BytesLib for bytes;

    function doubleKeccak256(bytes memory body) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(keccak256(body)));
    }

    function parseVMFromLogs(Vm.Log memory log) public pure returns (IDeltaswap.VM memory vm_) {
        uint256 index = 0;

        // emitterAddress
        vm_.emitterAddress = bytes32(log.topics[1]);

        // sequence
        vm_.sequence = log.data.toUint64(index + 32 - 8);
        index += 32;

        // nonce
        vm_.nonce = log.data.toUint32(index + 32 - 4);
        index += 32;

        // skip random bytes
        index += 32;

        // consistency level
        vm_.consistencyLevel = log.data.toUint8(index + 32 - 1);
        index += 32;

        // length of payload
        uint256 payloadLen = log.data.toUint256(index);
        index += 32;

        vm_.payload = log.data.slice(index, payloadLen);
        index += payloadLen;

        // trailing bytes (due to 32 byte slot overlap)
        index += log.data.length - index;

        require(index == log.data.length, "failed to parse deltaswap message");
    }

    /**
     * @notice Finds published Deltaswap events in forge logs
     * @param logs The forge Vm.log captured when recording events during test execution
     */
    function fetchDeltaswapMessageFromLog(Vm.Log[] memory logs)
        public
        pure
        returns (Vm.Log[] memory)
    {
        uint256 count = 0;
        for (uint256 i = 0; i < logs.length; i++) {
            if (
                logs[i].topics[0]
                    == keccak256("LogMessagePublished(address,uint64,uint32,bytes,uint8)")
            ) {
                count += 1;
            }
        }

        // create log array to save published messages
        Vm.Log[] memory published = new Vm.Log[](count);

        uint256 publishedIndex = 0;
        for (uint256 i = 0; i < logs.length; i++) {
            if (
                logs[i].topics[0]
                    == keccak256("LogMessagePublished(address,uint64,uint32,bytes,uint8)")
            ) {
                published[publishedIndex] = logs[i];
                publishedIndex += 1;
            }
        }

        return published;
    }

    /**
     * @notice Encodes Deltaswap message body into bytes
     * @param vm_ Deltaswap VM struct
     * @return encodedObservation Deltaswap message body encoded into bytes
     */
    function encodeObservation(IDeltaswap.VM memory vm_)
        public
        pure
        returns (bytes memory encodedObservation)
    {
        encodedObservation = abi.encodePacked(
            vm_.timestamp,
            vm_.nonce,
            vm_.emitterChainId,
            vm_.emitterAddress,
            vm_.sequence,
            vm_.consistencyLevel,
            vm_.payload
        );
    }

    /**
     * @notice Formats and signs a simulated Deltaswap message using the emitted log from calling `publishMessage`
     * @param log The forge Vm.log captured when recording events during test execution
     * @return signedMessage Formatted and signed Deltaswap message
     */
    function fetchSignedMessageFromLogs(
        Vm.Log memory log,
        uint16 emitterChainId,
        address emitterAddress
    ) public returns (bytes memory signedMessage) {
        // Parse deltaswap message from ethereum logs
        IDeltaswap.VM memory vm_ = parseVMFromLogs(log);

        // Set empty body values before computing the hash
        vm_.version = uint8(1);
        vm_.timestamp = uint32(block.timestamp);
        vm_.emitterChainId = emitterChainId;
        vm_.emitterAddress = bytes32(uint256(uint160(emitterAddress)));

        return encodeAndSignMessage(vm_);
    }

    /**
     * Functions that must be implemented by concrete deltaswap simulators.
     */

    /**
     * @notice Sets the message fee for a deltaswap message.
     */
    function setMessageFee(uint256 newFee) public virtual;

    /**
     * @notice Invalidates a VM. It must be executed before it is parsed and verified by the Deltaswap instance to work.
     */
    function invalidateVM(bytes memory message) public virtual;

    /**
     * @notice Signs and preformatted simulated Deltaswap message
     * @param vm_ The preformatted Deltaswap message
     * @return signedMessage Formatted and signed Deltaswap message
     */
    function encodeAndSignMessage(IDeltaswap.VM memory vm_)
        public
        virtual
        returns (bytes memory signedMessage);
}

/**
 * @title A Deltaswap Phylax Simulator
 * @notice This contract simulates signing Deltaswap messages emitted in a forge test.
 * This particular version doesn't sign any message but just exists to keep a standard interface for tests.
 * @dev This contract is meant to be used with the MockDeltaswap contract that validates any VM as long
 *   as its hash wasn't banned.
 */
contract FakeDeltaswapSimulator is DeltaswapSimulator {
    // Allow access to Deltaswap
    MockDeltaswap public deltaswap;

    /**
     * @param initDeltaswap address of the Deltaswap core contract for the mainnet chain being forked
     */
    constructor(MockDeltaswap initDeltaswap) {
        deltaswap = initDeltaswap;
    }

    function setMessageFee(uint256 newFee) public override {
        deltaswap.setMessageFee(newFee);
    }

    function invalidateVM(bytes memory message) public override {
        deltaswap.invalidateVM(message);
    }

    /**
     * @notice Signs and preformatted simulated Deltaswap message
     * @param vm_ The preformatted Deltaswap message
     * @return signedMessage Formatted and signed Deltaswap message
     */
    function encodeAndSignMessage(IDeltaswap.VM memory vm_)
        public
        view
        override
        returns (bytes memory signedMessage)
    {
        // Compute the hash of the body
        bytes memory body = encodeObservation(vm_);
        vm_.hash = doubleKeccak256(body);

        signedMessage = abi.encodePacked(
            vm_.version,
            deltaswap.getCurrentPhylaxSetIndex(),
            // length of signature array
            uint8(1),
            // phylax index
            uint8(0),
            // r sig argument
            bytes32(uint256(0)),
            // s sig argument
            bytes32(uint256(0)),
            // v sig argument (encodes public key recovery id, public key type and network of the signature)
            uint8(0),
            body
        );
    }
}

/**
 * @title A Deltaswap Phylax Simulator
 * @notice This contract simulates signing Deltaswap messages emitted in a forge test.
 * It overrides the Deltaswap phylax set to allow for signing messages with a single
 * private key on any EVM where Deltaswap core contracts are deployed.
 * @dev This contract is meant to be used when testing against a mainnet fork.
 */
contract SigningDeltaswapSimulator is DeltaswapSimulator {
    // Taken from forge-std/Script.sol
    address private constant VM_ADDRESS =
        address(bytes20(uint160(uint256(keccak256("hevm cheat code")))));
    Vm public constant vm = Vm(VM_ADDRESS);

    // Allow access to Deltaswap
    IDeltaswap public deltaswap;

    // Save the phylax PK to sign messages with
    uint256 private devnetPhylaxPK;

    /**
     * @param deltaswap_ address of the Deltaswap core contract for the mainnet chain being forked
     * @param devnetPhylax private key of the devnet Phylax
     */
    constructor(IDeltaswap deltaswap_, uint256 devnetPhylax) {
        deltaswap = deltaswap_;
        devnetPhylaxPK = devnetPhylax;
        overrideToDevnetPhylax(vm.addr(devnetPhylax));
    }

    function overrideToDevnetPhylax(address devnetPhylax) internal {
        {
            // Get slot for Phylax Set at the current index
            uint32 phylaxSetIndex = deltaswap.getCurrentPhylaxSetIndex();
            bytes32 phylaxSetSlot = keccak256(abi.encode(phylaxSetIndex, 2));

            // Overwrite all but first phylax set to zero address. This isn't
            // necessary, but just in case we inadvertently access these slots
            // for any reason.
            uint256 numPhylaxs = uint256(vm.load(address(deltaswap), phylaxSetSlot));
            for (uint256 i = 1; i < numPhylaxs;) {
                vm.store(
                    address(deltaswap),
                    bytes32(uint256(keccak256(abi.encodePacked(phylaxSetSlot))) + i),
                    bytes32(0)
                );
                unchecked {
                    i += 1;
                }
            }

            // Now overwrite the first phylax key with the devnet key specified
            // in the function argument.
            vm.store(
                address(deltaswap),
                bytes32(uint256(keccak256(abi.encodePacked(phylaxSetSlot))) + 0), // just explicit w/ index 0
                bytes32(uint256(uint160(devnetPhylax)))
            );

            // Change the length to 1 phylax
            vm.store(
                address(deltaswap),
                phylaxSetSlot,
                bytes32(uint256(1)) // length == 1
            );

            // Confirm phylax set override
            address[] memory phylaxs = deltaswap.getPhylaxSet(phylaxSetIndex).keys;
            require(phylaxs.length == 1, "phylaxs.length != 1");
            require(phylaxs[0] == devnetPhylax, "incorrect phylax set override");
        }
    }

    function setMessageFee(uint256 newFee) public override {
        bytes32 coreModule = 0x00000000000000000000000000000000000000000000000000000000436f7265;
        bytes memory message =
            abi.encodePacked(coreModule, uint8(3), uint16(deltaswap.chainId()), newFee);
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

        bytes memory signed = encodeAndSignMessage(preSignedMessage);
        deltaswap.submitSetMessageFee(signed);
    }

    function invalidateVM(bytes memory message) public pure override {
        // Don't do anything. Signatures are easily invalidated modifying the payload.
        // If it becomes necessary to prevent producing a good signature for this message, that can be done here.
    }

    /**
     * @notice Signs and preformatted simulated Deltaswap message
     * @param vm_ The preformatted Deltaswap message
     * @return signedMessage Formatted and signed Deltaswap message
     */
    function encodeAndSignMessage(IDeltaswap.VM memory vm_)
        public
        view
        override
        returns (bytes memory signedMessage)
    {
        // Compute the hash of the body
        bytes memory body = encodeObservation(vm_);
        vm_.hash = doubleKeccak256(body);

        // Sign the hash with the devnet phylax private key
        IDeltaswap.Signature[] memory sigs = new IDeltaswap.Signature[](1);
        (sigs[0].v, sigs[0].r, sigs[0].s) = vm.sign(devnetPhylaxPK, vm_.hash);
        sigs[0].phylaxIndex = 0;

        signedMessage = abi.encodePacked(
            vm_.version,
            deltaswap.getCurrentPhylaxSetIndex(),
            uint8(sigs.length),
            sigs[0].phylaxIndex,
            sigs[0].r,
            sigs[0].s,
            sigs[0].v - 27,
            body
        );
    }
}
