import { ethers, BigNumber } from "ethers";
import { ChainId, ChainName, CHAINS, Network } from "../../utils";
import { parseVaa } from "../../vaa";
import { getDeltaswapRelayer } from "../consts";
import {
  VaaKey,
  parseDeltaswapRelayerSend,
  parseEVMExecutionInfoV1,
} from "../structs";
import { vaaKeyToVaaKeyStruct, getDeliveryProvider } from "./helpers";
import { getSignedVAAWithRetry } from "../../rpc";

export async function resendRaw(
  signer: ethers.Signer,
  sourceChain: ChainName,
  targetChain: ChainName,
  environment: Network,
  vaaKey: VaaKey,
  newGasLimit: BigNumber | number,
  newReceiverValue: BigNumber | number,
  deliveryProviderAddress: string,
  overrides?: ethers.PayableOverrides,
  deltaswapRelayerAddress?: string
): Promise<ethers.providers.TransactionResponse> {
  const provider = signer.provider;

  if (!provider) throw Error("No provider on signer");

  const deltaswapRelayer = getDeltaswapRelayer(
    sourceChain,
    environment,
    signer,
    deltaswapRelayerAddress
  );

  return deltaswapRelayer.resendToEvm(
    vaaKeyToVaaKeyStruct(vaaKey),
    CHAINS[targetChain],
    newReceiverValue,
    newGasLimit,
    deliveryProviderAddress,
    overrides
  );
}

type ResendOptionalParams = {
  deltaswapRelayerAddress?: string;
};

export async function resend(
  signer: ethers.Signer,
  sourceChain: ChainName,
  targetChain: ChainName,
  environment: Network,
  vaaKey: VaaKey,
  newGasLimit: BigNumber | number,
  newReceiverValue: BigNumber | number,
  deliveryProviderAddress: string,
  deltaswapRPCs: string[],
  overrides: ethers.PayableOverrides,
  extraGrpcOpts = {},
  optionalParams?: ResendOptionalParams
): Promise<ethers.providers.TransactionResponse> {
  const targetChainId = CHAINS[targetChain];
  const originalVAA = await getSignedVAAWithRetry(
    deltaswapRPCs,
    vaaKey.chainId as ChainId,
    vaaKey.emitterAddress.toString("hex"),
    vaaKey.sequence.toBigInt().toString(),
    extraGrpcOpts,
    2000,
    4
  );

  if (!originalVAA.vaaBytes) throw Error("original VAA not found");

  const originalVAAparsed = parseDeltaswapRelayerSend(
    parseVaa(Buffer.from(originalVAA.vaaBytes)).payload
  );
  if (!originalVAAparsed) throw Error("original VAA not a valid delivery VAA.");

  const [originalExecutionInfo] = parseEVMExecutionInfoV1(
    originalVAAparsed.encodedExecutionInfo,
    0
  );
  const originalGasLimit = originalExecutionInfo.gasLimit;
  const originalRefund = originalExecutionInfo.targetChainRefundPerGasUnused;
  const originalReceiverValue = originalVAAparsed.requestedReceiverValue;
  const originalTargetChain = originalVAAparsed.targetChainId;

  if (originalTargetChain != targetChainId) {
    throw Error(
      `Target chain of original VAA (${originalTargetChain}) does not match target chain of resend (${targetChainId})`
    );
  }

  if (newReceiverValue < originalReceiverValue) {
    throw Error(
      `New receiver value too low. Minimum is ${originalReceiverValue.toString()}`
    );
  }

  if (newGasLimit < originalGasLimit) {
    throw Error(
      `New gas limit too low. Minimum is ${originalReceiverValue.toString()}`
    );
  }

  const deltaswapRelayer = getDeltaswapRelayer(
    sourceChain,
    environment,
    signer,
    optionalParams?.deltaswapRelayerAddress
  );

  const [deliveryPrice, refundPerUnitGas]: [BigNumber, BigNumber] =
    await deltaswapRelayer[
      "quoteEVMDeliveryPrice(uint16,uint256,uint256,address)"
    ](
      targetChainId,
      newReceiverValue || 0,
      newGasLimit,
      deliveryProviderAddress
    );
  const value = await (overrides?.value || 0);
  if (!deliveryPrice.eq(value)) {
    throw new Error(
      `Expected a payment of ${deliveryPrice.toString()} wei; received ${value.toString()} wei`
    );
  }

  if (refundPerUnitGas < originalRefund) {
    throw Error(
      `New refund per unit gas too low. Minimum is ${originalRefund.toString()}.`
    );
  }

  return resendRaw(
    signer,
    sourceChain,
    targetChain,
    environment,
    vaaKey,
    newGasLimit,
    newReceiverValue,
    deliveryProviderAddress,
    overrides,
    optionalParams?.deltaswapRelayerAddress
  );
}
