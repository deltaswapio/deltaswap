// contracts/Messages.sol
// SPDX-License-Identifier: Apache 2

pragma solidity ^0.8.0;
pragma experimental ABIEncoderV2;

import "./Getters.sol";
import "./Structs.sol";
import "./libraries/external/BytesLib.sol";


contract Messages is Getters {
    using BytesLib for bytes;

    /// @dev parseAndVerifyVM serves to parse an encodedVM and wholy validate it for consumption
    function parseAndVerifyVM(bytes calldata encodedVM) public view returns (Structs.VM memory vm, bool valid, string memory reason) {
        vm = parseVM(encodedVM);
        /// setting checkHash to false as we can trust the hash field in this case given that parseVM computes and then sets the hash field above
        (valid, reason) = verifyVMInternal(vm, false);
    }

   /**
    * @dev `verifyVM` serves to validate an arbitrary vm against a valid Phylax set
    *  - it aims to make sure the VM is for a known phylaxSet
    *  - it aims to ensure the phylaxSet is not expired
    *  - it aims to ensure the VM has reached quorum
    *  - it aims to verify the signatures provided against the phylaxSet
    *  - it aims to verify the hash field provided against the contents of the vm
    */
    function verifyVM(Structs.VM memory vm) public view returns (bool valid, string memory reason) {
        (valid, reason) = verifyVMInternal(vm, true);    
    }

    /**
    * @dev `verifyVMInternal` serves to validate an arbitrary vm against a valid Phylax set
    * if checkHash is set then the hash field of the vm is verified against the hash of its contents
    * in the case that the vm is securely parsed and the hash field can be trusted, checkHash can be set to false
    * as the check would be redundant
    */
    function verifyVMInternal(Structs.VM memory vm, bool checkHash) internal view returns (bool valid, string memory reason) {
        /// @dev Obtain the current phylaxSet for the phylaxSetIndex provided
        Structs.PhylaxSet memory phylaxSet = getPhylaxSet(vm.phylaxSetIndex);

        /**
         * Verify that the hash field in the vm matches with the hash of the contents of the vm if checkHash is set
         * WARNING: This hash check is critical to ensure that the vm.hash provided matches with the hash of the body.
         * Without this check, it would not be safe to call verifyVM on it's own as vm.hash can be a valid signed hash
         * but the body of the vm could be completely different from what was actually signed by the phylaxs
         */
        if(checkHash){
            bytes memory body = abi.encodePacked(
                vm.timestamp,
                vm.nonce,
                vm.emitterChainId,
                vm.emitterAddress,
                vm.sequence,
                vm.consistencyLevel,
                vm.payload
            );

            bytes32 vmHash = keccak256(abi.encodePacked(keccak256(body)));

            if(vmHash != vm.hash){
                return (false, "vm.hash doesn't match body");
            }
        }

       /**
        * @dev Checks whether the phylaxSet has zero keys
        * WARNING: This keys check is critical to ensure the phylaxSet has keys present AND to ensure
        * that phylaxSet key size doesn't fall to zero and negatively impact quorum assessment.  If phylaxSet
        * key length is 0 and vm.signatures length is 0, this could compromise the integrity of both vm and
        * signature verification.
        */
        if(phylaxSet.keys.length == 0){
            return (false, "invalid phylax set");
        }

        /// @dev Checks if VM phylax set index matches the current index (unless the current set is expired).
        if(vm.phylaxSetIndex != getCurrentPhylaxSetIndex() && phylaxSet.expirationTime < block.timestamp){
            return (false, "phylax set has expired");
        }

       /**
        * @dev We're using a fixed point number transformation with 1 decimal to deal with rounding.
        *   WARNING: This quorum check is critical to assessing whether we have enough Phylax signatures to validate a VM
        *   if making any changes to this, obtain additional peer review. If phylaxSet key length is 0 and
        *   vm.signatures length is 0, this could compromise the integrity of both vm and signature verification.
        */
        if (vm.signatures.length < quorum(phylaxSet.keys.length)){
            return (false, "no quorum");
        }

        /// @dev Verify the proposed vm.signatures against the phylaxSet
        (bool signaturesValid, string memory invalidReason) = verifySignatures(vm.hash, vm.signatures, phylaxSet);
        if(!signaturesValid){
            return (false, invalidReason);
        }

        /// If we are here, we've validated the VM is a valid multi-sig that matches the phylaxSet.
        return (true, "");
    }


    /**
     * @dev verifySignatures serves to validate arbitrary sigatures against an arbitrary phylaxSet
     *  - it intentionally does not solve for expectations within phylaxSet (you should use verifyVM if you need these protections)
     *  - it intentioanlly does not solve for quorum (you should use verifyVM if you need these protections)
     *  - it intentionally returns true when signatures is an empty set (you should use verifyVM if you need these protections)
     */
    function verifySignatures(bytes32 hash, Structs.Signature[] memory signatures, Structs.PhylaxSet memory phylaxSet) public pure returns (bool valid, string memory reason) {
        uint8 lastIndex = 0;
        uint256 phylaxCount = phylaxSet.keys.length;
        for (uint i = 0; i < signatures.length; i++) {
            Structs.Signature memory sig = signatures[i];
            address signatory = ecrecover(hash, sig.v, sig.r, sig.s);
            // ecrecover returns 0 for invalid signatures. We explicitly require valid signatures to avoid unexpected
            // behaviour due to the default storage slot value also being 0.
            require(signatory != address(0), "ecrecover failed with signature");

            /// Ensure that provided signature indices are ascending only
            require(i == 0 || sig.phylaxIndex > lastIndex, "signature indices must be ascending");
            lastIndex = sig.phylaxIndex;

            /// @dev Ensure that the provided signature index is within the
            /// bounds of the phylaxSet. This is implicitly checked by the array
            /// index operation below, so this check is technically redundant.
            /// However, reverting explicitly here ensures that a bug is not
            /// introduced accidentally later due to the nontrivial storage
            /// semantics of solidity.
            require(sig.phylaxIndex < phylaxCount, "phylax index out of bounds");

            /// Check to see if the signer of the signature does not match a specific Phylax key at the provided index
            if(signatory != phylaxSet.keys[sig.phylaxIndex]){
                return (false, "VM signature invalid");
            }
        }

        /// If we are here, we've validated that the provided signatures are valid for the provided phylaxSet
        return (true, "");
    }

    /**
     * @dev parseVM serves to parse an encodedVM into a vm struct
     *  - it intentionally performs no validation functions, it simply parses raw into a struct
     */
    function parseVM(bytes memory encodedVM) public pure virtual returns (Structs.VM memory vm) {
        uint index = 0;

        vm.version = encodedVM.toUint8(index);
        index += 1;
        // SECURITY: Note that currently the VM.version is not part of the hash 
        // and for reasons described below it cannot be made part of the hash. 
        // This means that this field's integrity is not protected and cannot be trusted. 
        // This is not a problem today since there is only one accepted version, but it 
        // could be a problem if we wanted to allow other versions in the future. 
        require(vm.version == 1, "VM version incompatible"); 

        vm.phylaxSetIndex = encodedVM.toUint32(index);
        index += 4;

        // Parse Signatures
        uint256 signersLen = encodedVM.toUint8(index);
        index += 1;
        vm.signatures = new Structs.Signature[](signersLen);
        for (uint i = 0; i < signersLen; i++) {
            vm.signatures[i].phylaxIndex = encodedVM.toUint8(index);
            index += 1;

            vm.signatures[i].r = encodedVM.toBytes32(index);
            index += 32;
            vm.signatures[i].s = encodedVM.toBytes32(index);
            index += 32;
            vm.signatures[i].v = encodedVM.toUint8(index) + 27;
            index += 1;
        }

        /*
        Hash the body

        SECURITY: Do not change the way the hash of a VM is computed! 
        Changing it could result into two different hashes for the same observation. 
        But xDapps rely on the hash of an observation for replay protection.
        */
        bytes memory body = encodedVM.slice(index, encodedVM.length - index);
        vm.hash = keccak256(abi.encodePacked(keccak256(body)));

        // Parse the body
        vm.timestamp = encodedVM.toUint32(index);
        index += 4;

        vm.nonce = encodedVM.toUint32(index);
        index += 4;

        vm.emitterChainId = encodedVM.toUint16(index);
        index += 2;

        vm.emitterAddress = encodedVM.toBytes32(index);
        index += 32;

        vm.sequence = encodedVM.toUint64(index);
        index += 8;

        vm.consistencyLevel = encodedVM.toUint8(index);
        index += 1;

        vm.payload = encodedVM.slice(index, encodedVM.length - index);
    }

    /**
     * @dev quorum serves solely to determine the number of signatures required to acheive quorum
     */
    function quorum(uint numPhylaxs) public pure virtual returns (uint numSignaturesRequiredForQuorum) {
        // The max number of phylaxs is 255
        require(numPhylaxs < 256, "too many phylaxs");
        return ((numPhylaxs * 2) / 3) + 1;
    }
}
