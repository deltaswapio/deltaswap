import {
  PublicKey,
  PublicKeyInitData,
  TransactionInstruction,
} from "@solana/web3.js";
import { createReadOnlyTokenBridgeProgramInterface } from "../program";
import { getPostMessageAccounts } from "../../deltaswap";
import {
  deriveTokenMetadataKey,
  deriveTokenBridgeConfigKey,
  deriveWrappedMetaKey,
} from "../accounts";

export function createAttestTokenInstruction(
  tokenBridgeProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  mint: PublicKeyInitData,
  message: PublicKeyInitData,
  nonce: number
): TransactionInstruction {
  const methods =
    createReadOnlyTokenBridgeProgramInterface(
      tokenBridgeProgramId
    ).methods.attestToken(nonce);

  // @ts-ignore
  return methods._ixFn(...methods._args, {
    accounts: getAttestTokenAccounts(
      tokenBridgeProgramId,
      deltaswapProgramId,
      payer,
      mint,
      message
    ) as any,
    signers: undefined,
    remainingAccounts: undefined,
    preInstructions: undefined,
    postInstructions: undefined,
  });
}

export interface AttestTokenAccounts {
  payer: PublicKey;
  config: PublicKey;
  mint: PublicKey;
  wrappedMeta: PublicKey;
  splMetadata: PublicKey;
 deltaswapBridge: PublicKey;
 deltaswapMessage: PublicKey;
  deltaswapEmitter: PublicKey;
  deltaswapSequence: PublicKey;
  deltaswapFeeCollector: PublicKey;
  clock: PublicKey;
  rent: PublicKey;
  systemProgram: PublicKey;
  deltaswapProgram: PublicKey;
}

export function getAttestTokenAccounts(
  tokenBridgeProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  mint: PublicKeyInitData,
  message: PublicKeyInitData
): AttestTokenAccounts {
  const {
    bridge:deltaswapBridge,
    emitter: deltaswapEmitter,
    sequence: deltaswapSequence,
    feeCollector: deltaswapFeeCollector,
    clock,
    rent,
    systemProgram,
  } = getPostMessageAccounts(
    deltaswapProgramId,
    payer,
    tokenBridgeProgramId,
    message
  );
  return {
    payer: new PublicKey(payer),
    config: deriveTokenBridgeConfigKey(tokenBridgeProgramId),
    mint: new PublicKey(mint),
    wrappedMeta: deriveWrappedMetaKey(tokenBridgeProgramId, mint),
    splMetadata: deriveTokenMetadataKey(mint),
   deltaswapBridge,
   deltaswapMessage: new PublicKey(message),
    deltaswapEmitter,
    deltaswapSequence,
    deltaswapFeeCollector,
    clock,
    rent,
    systemProgram,
    deltaswapProgram: new PublicKey(deltaswapProgramId),
  };
}
