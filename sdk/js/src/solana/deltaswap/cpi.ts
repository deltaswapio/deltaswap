import { PublicKey, PublicKeyInitData } from "@solana/web3.js";
import {
  deriveEmitterSequenceKey,
  deriveFeeCollectorKey,
  deriveDeltaswapEmitterKey,
  deriveDeltaswapBridgeDataKey,
  getEmitterKeys,
} from "./accounts";
import { getPostMessageAccounts } from "./instructions";

export interface DeltaswapDerivedAccounts {
  /**
   * seeds = ["Bridge"], seeds::program = deltaswapProgram
   */
 deltaswapBridge: PublicKey;
  /**
   * seeds = ["emitter"], seeds::program = cpiProgramId
   */
  deltaswapEmitter: PublicKey;
  /**
   * seeds = ["Sequence", deltaswapEmitter], seeds::program = deltaswapProgram
   */
  deltaswapSequence: PublicKey;
  /**
   * seeds = ["fee_collector"], seeds::program = deltaswapProgram
   */
  deltaswapFeeCollector: PublicKey;
}

/**
 * Generate Deltaswap PDAs.
 *
 * @param cpiProgramId
 * @param deltaswapProgramId
 * @returns
 */
export function getDeltaswapDerivedAccounts(
  cpiProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData
): DeltaswapDerivedAccounts {
  const { emitter: deltaswapEmitter, sequence: deltaswapSequence } =
    getEmitterKeys(cpiProgramId, deltaswapProgramId);
  return {
   deltaswapBridge: deriveDeltaswapBridgeDataKey(deltaswapProgramId),
    deltaswapEmitter,
    deltaswapSequence,
    deltaswapFeeCollector: deriveFeeCollectorKey(deltaswapProgramId),
  };
}

export interface PostMessageCpiAccounts extends DeltaswapDerivedAccounts {
  payer: PublicKey;
 deltaswapMessage: PublicKey;
  clock: PublicKey;
  rent: PublicKey;
  systemProgram: PublicKey;
}

/**
 * Generate accounts needed to perform `post_message` instruction
 * as cross-program invocation.
 *
 * @param cpiProgramId
 * @param deltaswapProgramId
 * @param payer
 * @param message
 * @returns
 */
export function getPostMessageCpiAccounts(
  cpiProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  message: PublicKeyInitData
): PostMessageCpiAccounts {
  const accounts = getPostMessageAccounts(
    deltaswapProgramId,
    payer,
    cpiProgramId,
    message
  );
  return {
    payer: accounts.payer,
   deltaswapBridge: accounts.bridge,
   deltaswapMessage: accounts.message,
    deltaswapEmitter: accounts.emitter,
    deltaswapSequence: accounts.sequence,
    deltaswapFeeCollector: accounts.feeCollector,
    clock: accounts.clock,
    rent: accounts.rent,
    systemProgram: accounts.systemProgram,
  };
}
