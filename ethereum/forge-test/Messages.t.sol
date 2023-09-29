// test/Messages.sol
// SPDX-License-Identifier: Apache 2

pragma solidity ^0.8.0;

import "../contracts/Messages.sol";
import "../contracts/Setters.sol";
import "../contracts/Structs.sol";
import "forge-std/Test.sol";

contract ExportedMessages is Messages, Setters {
    function storePhylaxSetPub(Structs.PhylaxSet memory set, uint32 index) public {
        return super.storePhylaxSet(set, index);
    }
}

contract TestMessages is Test {
  address constant testPhylaxPub = 0xbeFA429d57cD18b7F8A4d91A2da9AB4AF05d0FBe;

  // A valid VM with one signature from the testPhylaxPublic key
  bytes validVM = hex"01000000000100867b55fec41778414f0683e80a430b766b78801b7070f9198ded5e62f48ac7a44b379a6cf9920e42dbd06c5ebf5ec07a934a00a572aefc201e9f91c33ba766d900000003e800000001000b0000000000000000000000000000000000000000000000000000000000000eee00000000000005390faaaa";

  uint256 constant testPhylax = 93941733246223705020089879371323733820373732307041878556247502674739205313440;

  ExportedMessages messages;

  Structs.PhylaxSet phylaxSet;

  function setUp() public {
    messages = new ExportedMessages();

    // initialize phylax set with one phylax
    address[] memory keys = new address[](1);
    keys[0] = vm.addr(testPhylax);
    phylaxSet = Structs.PhylaxSet(keys, 0);
    require(messages.quorum(phylaxSet.keys.length) == 1, "Quorum should be 1");
  }

  function testQuorum() public {
    assertEq(messages.quorum(0), 1);
    assertEq(messages.quorum(1), 1);
    assertEq(messages.quorum(2), 2);
    assertEq(messages.quorum(3), 3);
    assertEq(messages.quorum(4), 3);
    assertEq(messages.quorum(5), 4);
    assertEq(messages.quorum(6), 5);
    assertEq(messages.quorum(7), 5);
    assertEq(messages.quorum(8), 6);
    assertEq(messages.quorum(9), 7);
    assertEq(messages.quorum(10), 7);
    assertEq(messages.quorum(11), 8);
    assertEq(messages.quorum(12), 9);
    assertEq(messages.quorum(19), 13);
    assertEq(messages.quorum(20), 14);
  }

  function testQuorumCanAlwaysBeReached(uint256 numPhylaxs) public {
    vm.assume(numPhylaxs > 0);

    if (numPhylaxs >= 256) {
      vm.expectRevert("too many phylaxs");
    }
    // test that quorums is never greater than the number of phylaxs
    assert(messages.quorum(numPhylaxs) <= numPhylaxs);
  }

  // This test ensures that submitting more signatures than expected will
  // trigger a "phylax index out of bounds" error.
  function testCannotVerifySignaturesWithOutOfBoundsSignature(bytes memory encoded) public {
    vm.assume(encoded.length > 0);

    bytes32 message = keccak256(encoded);

    // First generate a legitimate signature.
    Structs.Signature memory goodSignature = Structs.Signature(message, 0, 0, 0);
    (goodSignature.v, goodSignature.r, goodSignature.s) = vm.sign(testPhylax, message);
    assertEq(ecrecover(message, goodSignature.v, goodSignature.r, goodSignature.s), vm.addr(testPhylax));

    // Reuse legitimate signature above for the next signature. This will
    // bypass the "invalid signature" revert.
    Structs.Signature memory outOfBoundsSignature = goodSignature;
    outOfBoundsSignature.phylaxIndex = 1;

    // Attempt to verify signatures.
    Structs.Signature[] memory sigs = new Structs.Signature[](2);
    sigs[0] = goodSignature;
    sigs[1] = outOfBoundsSignature;

    vm.expectRevert("phylax index out of bounds");
    messages.verifySignatures(message, sigs, phylaxSet);
  }

  // This test ensures that submitting an invalid signature fails when
  // verifySignatures is called. Calling ecrecover should fail.
  function testCannotVerifySignaturesWithInvalidSignature(bytes memory encoded) public {
    vm.assume(encoded.length > 0);

    bytes32 message = keccak256(encoded);

    // Generate an invalid signature.
    Structs.Signature memory badSignature = Structs.Signature(message, 0, 0, 0);
    assertEq(ecrecover(message, badSignature.v, badSignature.r, badSignature.s), address(0));

    // Attempt to verify signatures.
    Structs.Signature[] memory sigs = new Structs.Signature[](2);
    sigs[0] = badSignature;

    vm.expectRevert("ecrecover failed with signature");
    messages.verifySignatures(message, sigs, phylaxSet);
  }

  function testVerifySignatures(bytes memory encoded) public {
    vm.assume(encoded.length > 0);

    bytes32 message = keccak256(encoded);

    // Generate legitimate signature.
    Structs.Signature memory goodSignature;
    (goodSignature.v, goodSignature.r, goodSignature.s) = vm.sign(testPhylax, message);
    assertEq(ecrecover(message, goodSignature.v, goodSignature.r, goodSignature.s), vm.addr(testPhylax));
    goodSignature.phylaxIndex = 0;

    // Attempt to verify signatures.
    Structs.Signature[] memory sigs = new Structs.Signature[](1);
    sigs[0] = goodSignature;

    (bool valid, string memory reason) = messages.verifySignatures(message, sigs, phylaxSet);
    assertEq(valid, true);
    assertEq(bytes(reason).length, 0);
  }

  // This test checks the possibility of getting a unsigned message verified through verifyVM
  function testHashMismatchedVMIsNotVerified() public {
    // Set the initial phylax set
    address[] memory initialPhylaxs = new address[](1);
    initialPhylaxs[0] = testPhylaxPub;

    // Create a phylax set
    Structs.PhylaxSet memory initialPhylaxSet = Structs.PhylaxSet({
      keys: initialPhylaxs,
      expirationTime: 0
    });

    messages.storePhylaxSetPub(initialPhylaxSet, uint32(0));

    // Confirm that the test VM is valid
    (Structs.VM memory parsedValidVm, bool valid, string memory reason) = messages.parseAndVerifyVM(validVM);
    require(valid, reason);
    assertEq(valid, true);
    assertEq(reason, "");

    // Manipulate the payload of the vm
    Structs.VM memory invalidVm = parsedValidVm;
    invalidVm.payload = abi.encodePacked(
        parsedValidVm.payload,
        "malicious bytes in payload"
    );

    // Confirm that the verifyVM fails on invalid VM
    (valid, reason) = messages.verifyVM(invalidVm);
    assertEq(valid, false);
    assertEq(reason, "vm.hash doesn't match body");
  }
}
