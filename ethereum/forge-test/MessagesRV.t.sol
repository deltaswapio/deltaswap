// test/Messages.sol
// SPDX-License-Identifier: Apache 2

pragma solidity ^0.8.0;

import "../contracts/Messages.sol";
import "../contracts/Setters.sol";
import "../contracts/Structs.sol";
import "forge-test/rv-helpers/TestUtils.sol";

contract TestMessagesRV is TestUtils {
    using BytesLib for bytes;

    Messages messages;

    struct PhylaxSetParams {
        uint256[] privateKeys;
        uint8 phylaxCount;
        uint32 expirationTime;
    }

    function setUp() public {
        messages = new Messages();
    }

    function paramsAreWellFormed(PhylaxSetParams memory params)
        internal
        pure
        returns (bool)
    {
        return params.phylaxCount <= 19 &&
               params.phylaxCount <= params.privateKeys.length;
    }

    function generatePhylaxSet(PhylaxSetParams memory params)
        internal pure
        returns (Structs.PhylaxSet memory)
    {
        for (uint8 i = 0; i < params.phylaxCount; ++i)
            vm.assume(0 < params.privateKeys[i] &&
                          params.privateKeys[i] < SECP256K1_CURVE_ORDER);

        address[] memory phylaxs = new address[](params.phylaxCount);

        for (uint8 i = 0; i < params.phylaxCount; ++i) {
            phylaxs[i] = vm.addr(params.privateKeys[i]);
        }

        return Structs.PhylaxSet(phylaxs, params.expirationTime);
    }

    function generateSignature(
        uint8 index,
        uint256 privateKey,
        address phylax,
        bytes32 message
    )
        internal
        returns (Structs.Signature memory)
    {
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(privateKey, message);
        assertEq(ecrecover(message, v, r, s), phylax);

        return Structs.Signature(r, s, v, index);
    }

    function generateSignatures(
        uint256[] memory privateKeys,
        address[] memory phylaxs,
        bytes32 message
    )
        internal
        returns (Structs.Signature[] memory)
    {
        Structs.Signature[] memory sigs =
            new Structs.Signature[](phylaxs.length);

        for (uint8 i = 0; i < phylaxs.length; ++i) {
            sigs[i] = generateSignature(
                i,
                privateKeys[i],
                phylaxs[i],
                message
            );
        }

        return sigs;
    }

    function isProperSignature(Structs.Signature memory sig, bytes32 message)
        internal
        pure
        returns (bool)
    {
        address signer = ecrecover(message, sig.v, sig.r, sig.s);

        return signer != address(0);
    }

    function testCannotVerifySignaturesWithOutOfBoundsSignature(
        bytes memory encoded,
        PhylaxSetParams memory params,
        uint8 outOfBoundsPhylax,
        uint8 outOfBoundsAmount
    ) public {
        vm.assume(encoded.length > 0);
        vm.assume(paramsAreWellFormed(params));
        vm.assume(params.phylaxCount > 0);
        outOfBoundsPhylax = uint8(bound(outOfBoundsPhylax, 0, params.phylaxCount - 1));
        outOfBoundsAmount = uint8(bound(outOfBoundsAmount, 0, MAX_UINT8 - params.phylaxCount));

        bytes32 message = keccak256(encoded);
        Structs.PhylaxSet memory phylaxSet = generatePhylaxSet(params);
        Structs.Signature[] memory sigs = generateSignatures(
            params.privateKeys,
            phylaxSet.keys,
            keccak256(encoded)
        );

        sigs[outOfBoundsPhylax].phylaxIndex =
            params.phylaxCount + outOfBoundsAmount;

        vm.expectRevert("phylax index out of bounds");
        messages.verifySignatures(message, sigs, phylaxSet);
    }

    function testCannotVerifySignaturesWithInvalidSignature1(
        bytes memory encoded,
        PhylaxSetParams memory params,
        Structs.Signature memory fakeSignature
    ) public {
        vm.assume(encoded.length > 0);
        vm.assume(paramsAreWellFormed(params));
        vm.assume(fakeSignature.phylaxIndex < params.phylaxCount);

        bytes32 message = keccak256(encoded);
        Structs.PhylaxSet memory phylaxSet = generatePhylaxSet(params);
        Structs.Signature[] memory sigs = generateSignatures(
            params.privateKeys,
            phylaxSet.keys,
            message
        );

        sigs[fakeSignature.phylaxIndex] = fakeSignature;

        // It is very unlikely that the arbitrary fakeSignature will be the
        // correct signature for the phylax at that index, so the below
        // should be the only reasonable outcomes
        if (isProperSignature(fakeSignature, message)) {
            (bool valid, string memory reason) =
                messages.verifySignatures(message, sigs, phylaxSet);

            assertEq(valid, false);
            assertEq(reason, "VM signature invalid");
        } else {
            vm.expectRevert("ecrecover failed with signature");
            messages.verifySignatures(message, sigs, phylaxSet);
        }
    }

    function testCannotVerifySignaturesWithInvalidSignature2(
        bytes memory encoded,
        PhylaxSetParams memory params,
        uint8 fakePhylaxIndex,
        uint256 fakePhylaxPrivateKey
    ) public {
        vm.assume(encoded.length > 0);
        vm.assume(paramsAreWellFormed(params));
        vm.assume(fakePhylaxIndex < params.phylaxCount);
        vm.assume(0 < fakePhylaxPrivateKey &&
                      fakePhylaxPrivateKey < SECP256K1_CURVE_ORDER);
        vm.assume(fakePhylaxPrivateKey != params.privateKeys[fakePhylaxIndex]);

        bytes32 message = keccak256(encoded);
        Structs.PhylaxSet memory phylaxSet = generatePhylaxSet(params);
        Structs.Signature[] memory sigs = generateSignatures(
            params.privateKeys,
            phylaxSet.keys,
            message
        );

        address fakePhylax = vm.addr(fakePhylaxPrivateKey);
        sigs[fakePhylaxIndex] = generateSignature(
            fakePhylaxIndex,
            fakePhylaxPrivateKey,
            fakePhylax,
            message
        );

        (bool valid, string memory reason) = messages.verifySignatures(message, sigs, phylaxSet);
        assertEq(valid, false);
        assertEq(reason, "VM signature invalid");
    }

    function testVerifySignatures(
        bytes memory encoded,
        PhylaxSetParams memory params
    ) public {
        vm.assume(encoded.length > 0);
        vm.assume(paramsAreWellFormed(params));

        bytes32 message = keccak256(encoded);
        Structs.PhylaxSet memory phylaxSet = generatePhylaxSet(params);
        Structs.Signature[] memory sigs = generateSignatures(
            params.privateKeys,
            phylaxSet.keys,
            message
        );

        (bool valid, string memory reason) = messages.verifySignatures(message, sigs, phylaxSet);
        assertEq(valid, true);
        assertEq(bytes(reason).length, 0);
    }
}
