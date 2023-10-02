// test/Messages.sol
// SPDX-License-Identifier: Apache 2

pragma solidity ^0.8.0;

import "../contracts/Implementation.sol";
import "../contracts/Setup.sol";
import "../contracts/Deltaswap.sol";
import "../contracts/interfaces/IDeltaswap.sol";
import "forge-test/rv-helpers/TestUtils.sol";

contract TestSetup is TestUtils {

    Deltaswap proxy;
    Implementation impl;
    Setup setup;
    Setup proxiedSetup;
    IDeltaswap proxied;

    uint256 constant testPhylax = 93941733246223705020089879371323733820373732307041878556247502674739205313440;
    bytes32 constant governanceContract = 0x0000000000000000000000000000000000000000000000000000000000000004;

    function setUp() public {
        // Deploy setup
        setup = new Setup();
        // Deploy implementation contract
        impl = new Implementation();
        // Deploy proxy
        proxy = new Deltaswap(address(setup), bytes(""));

        address[] memory keys = new address[](1);
        keys[0] = 0xbeFA429d57cD18b7F8A4d91A2da9AB4AF05d0FBe; // vm.addr(testPhylax)

        //proxied setup
        proxiedSetup = Setup(address(proxy));

        vm.chainId(1);
        proxiedSetup.setup({
            implementation: address(impl),
            initialPhylaxs: keys,
            chainId: 2,
            governanceChainId: 1,
            governanceContract: governanceContract,
            evmChainId: 1
        });

        proxied = IDeltaswap(address(proxy));
    }

    function testInitialize_after_setup_revert(bytes32 storageSlot, address alice)
        public
        unchangedStorage(address(proxied), storageSlot)
    {
        vm.prank(alice);
        vm.expectRevert("already initialized");
        proxied.initialize();
    }

    function testInitialize_after_setup_revert_KEVM(bytes32 storageSlot, address alice)
        public
    {
        kevm.infiniteGas();
        testInitialize_after_setup_revert(storageSlot, alice);
    }

    function testSetup_after_setup_revert(
        bytes32 storageSlot,
        address alice,
        address implementation,
        address initialPhylax,
        uint16 chainId,
        uint16 governanceChainId,
        bytes32 govContract,
        uint256 evmChainId)
        public
        unchangedStorage(address(proxied), storageSlot)
    {
        address[] memory keys = new address[](1);
        keys[0] = initialPhylax;

        vm.prank(alice);
        vm.expectRevert("unsupported");
        proxiedSetup.setup({
            implementation: implementation,
            initialPhylaxs: keys,
            chainId: chainId,
            governanceChainId: governanceChainId,
            governanceContract: govContract,
            evmChainId: evmChainId
        });
    }

    function testSetup_after_setup_revert_KEVM(
        bytes32 storageSlot,
        address alice,
        address implementation,
        address initialPhylax,
        uint16 chainId,
        uint16 governanceChainId,
        bytes32 govContract,
        uint256 evmChainId)
        public
    {
        kevm.infiniteGas();
        testSetup_after_setup_revert(
            storageSlot,
            alice,
            implementation,
            initialPhylax,
            chainId,
            governanceChainId,
            govContract,
            evmChainId
        );
    }
}
