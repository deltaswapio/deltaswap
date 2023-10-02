// contracts/mock/MockBatchedVAASender.sol
// SPDX-License-Identifier: Apache 2

pragma solidity ^0.8.0;

import "../libraries/external/BytesLib.sol";
import "../interfaces/IDeltaswap.sol";

contract MockBatchedVAASender {
    using BytesLib for bytes;

    address deltaswapCoreAddress;

    function sendMultipleMessages(
        uint32 nonce,
        bytes memory payload,
        uint8 consistencyLevel
    )
        public
        payable
        returns (
            uint64 messageSequence0,
            uint64 messageSequence1,
            uint64 messageSequence2
        )
    {
        messageSequence0 = deltaswapCore().publishMessage{value: msg.value}(
            nonce,
            payload,
            consistencyLevel
        );

        messageSequence1 = deltaswapCore().publishMessage{value: msg.value}(
            nonce,
            payload,
            consistencyLevel
        );

        messageSequence2 = deltaswapCore().publishMessage{value: msg.value}(
            nonce,
            payload,
            consistencyLevel
        );
    }

    function deltaswapCore() private view returns (IDeltaswap) {
        return IDeltaswap(deltaswapCoreAddress);
    }

    function setup(address _deltaswapCore) public {
        deltaswapCoreAddress = _deltaswapCore;
    }
}
