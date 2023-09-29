import { ChainId, ChainName, getGovernorIsVAAEnqueued } from "..";
import { coalesceChainId } from "../utils";

export async function getGovernorIsVAAEnqueuedWithRetry(
  hosts: string[],
  emitterChain: ChainId | ChainName,
  emitterAddress: string,
  sequence: string,
  extraGrpcOpts = {},
  retryTimeout = 1000,
  retryAttempts?: number
) {
  let currentDeltaswapRpcHost = -1;
  const getNextRpcHost = () => ++currentDeltaswapRpcHost % hosts.length;
  let result;
  let attempts = 0;
  while (!result) {
    attempts++;
    await new Promise((resolve) => setTimeout(resolve, retryTimeout));
    try {
      result = await getGovernorIsVAAEnqueued(
        hosts[getNextRpcHost()],
        coalesceChainId(emitterChain),
        emitterAddress,
        sequence,
        extraGrpcOpts
      );
    } catch (e) {
      if (retryAttempts !== undefined && attempts > retryAttempts) {
        throw e;
      }
    }
  }
  return result;
}
