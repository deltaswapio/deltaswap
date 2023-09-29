import {
  ChainId,
  CHAIN_ID_TO_NAME,
  ChainName,
  isChain,
  CONTRACTS,
  CHAINS,
  tryNativeToHexString,
  Network,
  ethers_contracts,
} from "../..";
import { BigNumber, ethers } from "ethers";
import { getDeltaswapRelayerAddress } from "../consts";
import {
  RelayerPayloadId,
  DeliveryInstruction,
  RefundStatus,
  parseEVMExecutionInfoV1,
  DeliveryOverrideArgs,
  KeyType,
  parseVaaKey,
  parseCCTPKey,
} from "../structs";
import {
  getDefaultProvider,
  printChain,
  getDeltaswapRelayerLog,
  parseDeltaswapLog,
  getBlockRange,
  getDeltaswapRelayerInfoBySourceSequence,
} from "./helpers";
import { DeliveryInfo } from "./deliver";

export type InfoRequestParams = {
  environment?: Network;
  sourceChainProvider?: ethers.providers.Provider;
  targetChainProviders?: Map<ChainName, ethers.providers.Provider>;
  targetChainBlockRanges?: Map<
    ChainName,
    [ethers.providers.BlockTag, ethers.providers.BlockTag]
  >;
  deltaswapRelayerWhMessageIndex?: number;
  deltaswapRelayerAddresses?: Map<ChainName, string>;
};

export type GetPriceOptParams = {
  environment?: Network;
  receiverValue?: ethers.BigNumberish;
  deltaswapRelayerAddress?: string;
  deliveryProviderAddress?: string;
  sourceChainProvider?: ethers.providers.Provider;
};

export async function getPriceAndRefundInfo(
  sourceChain: ChainName,
  targetChain: ChainName,
  gasAmount: ethers.BigNumberish,
  optionalParams?: GetPriceOptParams
): Promise<[ethers.BigNumber, ethers.BigNumber]> {
  const environment = optionalParams?.environment || "MAINNET";
  const sourceChainProvider =
    optionalParams?.sourceChainProvider ||
    getDefaultProvider(environment, sourceChain);
  if (!sourceChainProvider)
    throw Error(
      "No default RPC for this chain; pass in your own provider (as sourceChainProvider)"
    );
  const deltaswapRelayerAddress =
    optionalParams?.deltaswapRelayerAddress ||
    getDeltaswapRelayerAddress(sourceChain, environment);
  const sourceDeltaswapRelayer =
    ethers_contracts.IDeltaswapRelayer__factory.connect(
      deltaswapRelayerAddress,
      sourceChainProvider
    );
  const deliveryProviderAddress =
    optionalParams?.deliveryProviderAddress ||
    (await sourceDeltaswapRelayer.getDefaultDeliveryProvider());
  const targetChainId = CHAINS[targetChain];
  const priceAndRefundInfo = await sourceDeltaswapRelayer[
    "quoteEVMDeliveryPrice(uint16,uint256,uint256,address)"
  ](
    targetChainId,
    optionalParams?.receiverValue || 0,
    gasAmount,
    deliveryProviderAddress
  );
  return priceAndRefundInfo;
}

export async function getPrice(
  sourceChain: ChainName,
  targetChain: ChainName,
  gasAmount: ethers.BigNumberish,
  optionalParams?: GetPriceOptParams
): Promise<ethers.BigNumber> {
  const priceAndRefundInfo = await getPriceAndRefundInfo(
    sourceChain,
    targetChain,
    gasAmount,
    optionalParams
  );
  return priceAndRefundInfo[0];
}

export async function getDeltaswapRelayerInfo(
  sourceChain: ChainName,
  sourceTransaction: string,
  infoRequest?: InfoRequestParams
): Promise<DeliveryInfo> {
  const environment = infoRequest?.environment || "MAINNET";
  const sourceChainProvider =
    infoRequest?.sourceChainProvider ||
    getDefaultProvider(environment, sourceChain);
  if (!sourceChainProvider)
    throw Error(
      "No default RPC for this chain; pass in your own provider (as sourceChainProvider)"
    );
  const receipt = await sourceChainProvider.getTransactionReceipt(
    sourceTransaction
  );
  if (!receipt) throw Error("Transaction has not been mined");
  const bridgeAddress = CONTRACTS[environment][sourceChain].core;
  const deltaswapRelayerAddress =
    infoRequest?.deltaswapRelayerAddresses?.get(sourceChain) ||
    getDeltaswapRelayerAddress(sourceChain, environment);
  if (!bridgeAddress || !deltaswapRelayerAddress) {
    throw Error(
      `Invalid chain ID or network: Chain ${sourceChain}, ${environment}`
    );
  }
  const deliveryLog = getDeltaswapRelayerLog(
    receipt,
    bridgeAddress,
    tryNativeToHexString(deltaswapRelayerAddress, "ethereum"),
    infoRequest?.deltaswapRelayerWhMessageIndex
      ? infoRequest.deltaswapRelayerWhMessageIndex
      : 0
  );

  const { type, parsed } = parseDeltaswapLog(deliveryLog.log);

  const instruction = parsed as DeliveryInstruction;

  const targetChainId = instruction.targetChainId as ChainId;
  if (!isChain(targetChainId)) throw Error(`Invalid Chain: ${targetChainId}`);
  const targetChain = CHAIN_ID_TO_NAME[targetChainId];
  const targetChainProvider =
    infoRequest?.targetChainProviders?.get(targetChain) ||
    getDefaultProvider(environment, targetChain);

  if (!targetChainProvider) {
    throw Error(
      "No default RPC for this chain; pass in your own provider (as targetChainProvider)"
    );
  }
  const [blockStartNumber, blockEndNumber] =
    infoRequest?.targetChainBlockRanges?.get(targetChain) ||
    getBlockRange(targetChainProvider);

  const targetChainStatus = await getDeltaswapRelayerInfoBySourceSequence(
    environment,
    targetChain,
    targetChainProvider,
    sourceChain,
    BigNumber.from(deliveryLog.sequence),
    blockStartNumber,
    blockEndNumber,
    infoRequest?.deltaswapRelayerAddresses?.get(targetChain) ||
      getDeltaswapRelayerAddress(targetChain, environment)
  );

  return {
    type: RelayerPayloadId.Delivery,
    sourceChain: sourceChain,
    sourceTransactionHash: sourceTransaction,
    sourceDeliverySequenceNumber: BigNumber.from(
      deliveryLog.sequence
    ).toNumber(),
    deliveryInstruction: instruction,
    targetChainStatus,
  };
}

export function printDeltaswapRelayerInfo(info: DeliveryInfo) {
  console.log(stringifyDeltaswapRelayerInfo(info));
}

export function stringifyDeltaswapRelayerInfo(
  info: DeliveryInfo,
  excludeSourceInformation?: boolean,
  overrides?: DeliveryOverrideArgs
): string {
  let stringifiedInfo = "";
  if (
    info.type == RelayerPayloadId.Delivery &&
    info.deliveryInstruction.targetAddress.toString("hex") !==
      "0000000000000000000000000000000000000000000000000000000000000000"
  ) {
    if (!excludeSourceInformation) {
      stringifiedInfo += `Found delivery request in transaction ${
        info.sourceTransactionHash
      } on ${
        info.sourceChain
      }\nfrom sender ${info.deliveryInstruction.senderAddress.toString(
        "hex"
      )} from ${info.sourceChain} with delivery sequence number ${
        info.sourceDeliverySequenceNumber
      }\n`;
    } else {
      stringifiedInfo += `Found delivery request from sender ${info.deliveryInstruction.senderAddress.toString(
        "hex"
      )}\n`;
    }
    const numMsgs = info.deliveryInstruction.messageKeys.length;

    const payload = info.deliveryInstruction.payload.toString("hex");
    if (payload.length > 0) {
      stringifiedInfo += `\nPayload to be relayed (as hex string): 0x${payload}`;
    }
    if (numMsgs > 0) {
      stringifiedInfo += `\nThe following ${numMsgs} deltaswap messages (VAAs) were ${
        payload.length > 0 ? "also " : ""
      }requested to be relayed:\n`;
      stringifiedInfo += info.deliveryInstruction.messageKeys
        .map((msgKey, i) => {
          let result = "";
          if (msgKey.keyType == KeyType.VAA) {
            const vaaKey = parseVaaKey(msgKey.key);
            result += `(VAA ${i}): `;
            result += `Message from ${
              vaaKey.chainId ? printChain(vaaKey.chainId) : ""
            }, with emitter address ${vaaKey.emitterAddress?.toString(
              "hex"
            )} and sequence number ${vaaKey.sequence}`;
          } else if (msgKey.keyType == KeyType.CCTP) {
            const cctpKey = parseCCTPKey(msgKey.key);
            result += `(CCTP ${i}): `;
            result += `Transfer from cctp domain ${printChain(cctpKey.domain)}`;
            result += `, with nonce ${cctpKey.nonce}`;
          } else {
            result += `(Unknown key type${i}): ${msgKey.keyType}`;
          }
          return result;
        })
        .join(",\n");
    }
    if (payload.length == 0 && numMsgs == 0) {
      stringifiedInfo += `\nAn empty payload was requested to be sent`;
    }

    const instruction = info.deliveryInstruction;
    if (overrides) {
      instruction.requestedReceiverValue = overrides.newReceiverValue;
      instruction.encodedExecutionInfo = overrides.newExecutionInfo;
    }

    const targetChainName =
      CHAIN_ID_TO_NAME[instruction.targetChainId as ChainId];
    stringifiedInfo += `${
      numMsgs == 0
        ? payload.length == 0
          ? ""
          : "\n\nPayload was requested to be relayed"
        : "\n\nThese were requested to be sent"
    } to 0x${instruction.targetAddress.toString("hex")} on ${printChain(
      instruction.targetChainId
    )}\n`;
    const totalReceiverValue = instruction.requestedReceiverValue.add(
      instruction.extraReceiverValue
    );
    stringifiedInfo += totalReceiverValue.gt(0)
      ? `Amount to pass into target address: ${ethers.utils.formatEther(
          totalReceiverValue
        )} of ${targetChainName} currency ${
          instruction.extraReceiverValue.gt(0)
            ? `\n${ethers.utils.formatEther(
                instruction.requestedReceiverValue
              )} requested, ${ethers.utils.formatEther(
                instruction.extraReceiverValue
              )} additionally paid for`
            : ""
        }\n`
      : ``;
    const [executionInfo] = parseEVMExecutionInfoV1(
      instruction.encodedExecutionInfo,
      0
    );
    stringifiedInfo += `Gas limit: ${executionInfo.gasLimit} ${targetChainName} gas\n`;

    const refundAddressChosen =
      instruction.refundAddress !== instruction.refundDeliveryProvider;
    if (refundAddressChosen) {
      stringifiedInfo += `Refund rate: ${ethers.utils.formatEther(
        executionInfo.targetChainRefundPerGasUnused
      )} of ${targetChainName} currency per unit of gas unused\n`;
      stringifiedInfo += `Refund address: ${instruction.refundAddress.toString(
        "hex"
      )}\n`;
    }
    stringifiedInfo += `\n`;
    stringifiedInfo += info.targetChainStatus.events

      .map(
        (e, i) =>
          `Delivery attempt ${i + 1}: ${
            e.transactionHash
              ? ` ${targetChainName} transaction hash: ${e.transactionHash}`
              : ""
          }\nStatus: ${e.status}\n${
            e.revertString
              ? `Failure reason: ${
                  e.gasUsed.eq(executionInfo.gasLimit)
                    ? "Gas limit hit"
                    : e.revertString
                }\n`
              : ""
          }Gas used: ${e.gasUsed.toString()}\nTransaction fee used: ${ethers.utils.formatEther(
            executionInfo.targetChainRefundPerGasUnused.mul(e.gasUsed)
          )} of ${targetChainName} currency\n${`Refund amount: ${ethers.utils.formatEther(
            executionInfo.targetChainRefundPerGasUnused.mul(
              executionInfo.gasLimit.sub(e.gasUsed)
            )
          )} of ${targetChainName} currency \nRefund status: ${
            e.refundStatus
          }\n`}`
      )
      .join("\n");
  } else if (
    info.type == RelayerPayloadId.Delivery &&
    info.deliveryInstruction.targetAddress.toString("hex") ===
      "0000000000000000000000000000000000000000000000000000000000000000"
  ) {
    stringifiedInfo += `Found delivery request in transaction ${info.sourceTransactionHash} on ${info.sourceChain}\n`;

    const instruction = info.deliveryInstruction;
    const targetChainName =
      CHAIN_ID_TO_NAME[instruction.targetChainId as ChainId];

    stringifiedInfo += `\nA refund of ${ethers.utils.formatEther(
      instruction.extraReceiverValue
    )} ${targetChainName} currency was requested to be sent to ${targetChainName}, address 0x${info.deliveryInstruction.refundAddress.toString(
      "hex"
    )}\n\n`;

    stringifiedInfo += info.targetChainStatus.events

      .map(
        (e, i) =>
          `Delivery attempt ${i + 1}: ${
            e.transactionHash
              ? ` ${targetChainName} transaction hash: ${e.transactionHash}`
              : ""
          }\nStatus: ${
            e.refundStatus == RefundStatus.RefundSent
              ? "Refund Successful"
              : "Refund Failed"
          }`
      )
      .join("\n");
  }

  return stringifiedInfo;
}
