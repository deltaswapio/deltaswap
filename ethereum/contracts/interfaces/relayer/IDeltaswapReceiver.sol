// SPDX-License-Identifier: Apache 2

pragma solidity ^0.8.0;

/**
 * @notice Interface for a contract which can receive Deltaswap messages.
 */
interface IDeltaswapReceiver {
    /**
     * @notice When a `send` is performed with this contract as the target, this function will be
     *     invoked by the DeltaswapRelayer contract
     *
     * NOTE: This function should be restricted such that only the Deltaswap Relayer contract can call it.
     *
     * We also recommend that this function checks that `sourceChain` and `sourceAddress` are indeed who
     *       you expect to have requested the calling of `send` on the source chain
     *
     * The invocation of this function corresponding to the `send` request will have msg.value equal
     *   to the receiverValue specified in the send request.
     *
     * If the invocation of this function reverts or exceeds the gas limit
     *   specified by the send requester, this delivery will result in a `ReceiverFailure`.
     *
     * @param payload - an arbitrary message which was included in the delivery by the
     *     requester. This message's signature will already have been verified (as long as msg.sender is the Deltaswap Relayer contract)
     * @param additionalMessages - Additional messages which were requested to be included in this delivery.
     *      Note: There are no contract-level guarantees that the messages in this array are what was requested
     *      so **you should verify any sensitive information given here!**
     *
     *      For example, if a 'VaaKey' was specified on the source chain, then MAKE SURE the corresponding message here
     *      has valid signatures (by calling `parseAndVerifyVM(message)` on the Deltaswap core contract)
     *
     *      This field can be used to perform and relay TokenBridge or CCTP transfers, and there are example
     *      usages of this at
     *         https://github.com/deltaswap-foundation/hello-token
     *         https://github.com/deltaswap-foundation/hello-cctp
     *
     * @param sourceAddress - the (deltaswap format) address on the sending chain which requested
     *     this delivery.
     * @param sourceChain - the deltaswap chain ID where this delivery was requested.
     * @param deliveryHash - the VAA hash of the deliveryVAA.
     *
     */
    function receiveDeltaswapMessages(
        bytes memory payload,
        bytes[] memory additionalMessages,
        bytes32 sourceAddress,
        uint16 sourceChain,
        bytes32 deliveryHash
    ) external payable;
}
