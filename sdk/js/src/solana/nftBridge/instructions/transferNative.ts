import { TOKEN_PROGRAM_ID } from "@solana/spl-token";
import {
  PublicKey,
  PublicKeyInitData,
  TransactionInstruction,
} from "@solana/web3.js";
import { TOKEN_METADATA_PROGRAM_ID, deriveTokenMetadataKey } from "../../utils";
import { getPostMessageAccounts } from "../../deltaswap";
import {
  deriveAuthoritySignerKey,
  deriveCustodyKey,
  deriveCustodySignerKey,
  deriveNftBridgeConfigKey,
} from "../accounts";
import { createReadOnlyNftBridgeProgramInterface } from "../program";

export function createTransferNativeInstruction(
  nftBridgeProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  message: PublicKeyInitData,
  from: PublicKeyInitData,
  mint: PublicKeyInitData,
  nonce: number,
  targetAddress: Buffer | Uint8Array,
  targetChain: number
): TransactionInstruction {
  const methods = createReadOnlyNftBridgeProgramInterface(
    nftBridgeProgramId
  ).methods.transferNative(
    nonce,
    Buffer.from(targetAddress) as any,
    targetChain
  );

  // @ts-ignore
  return methods._ixFn(...methods._args, {
    accounts: getTransferNativeAccounts(
      nftBridgeProgramId,
      deltaswapProgramId,
      payer,
      message,
      from,
      mint
    ) as any,
    signers: undefined,
    remainingAccounts: undefined,
    preInstructions: undefined,
    postInstructions: undefined,
  });
}

export interface TransferNativeAccounts {
  payer: PublicKey;
  config: PublicKey;
  from: PublicKey;
  mint: PublicKey;
  splMetadata: PublicKey;
  custody: PublicKey;
  authoritySigner: PublicKey;
  custodySigner: PublicKey;
 deltaswapBridge: PublicKey;
 deltaswapMessage: PublicKey;
  deltaswapEmitter: PublicKey;
  deltaswapSequence: PublicKey;
  deltaswapFeeCollector: PublicKey;
  clock: PublicKey;
  rent: PublicKey;
  systemProgram: PublicKey;
  tokenProgram: PublicKey;
  splMetadataProgram: PublicKey;
  deltaswapProgram: PublicKey;
}

export function getTransferNativeAccounts(
  nftBridgeProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  message: PublicKeyInitData,
  from: PublicKeyInitData,
  mint: PublicKeyInitData
): TransferNativeAccounts {
  const {
    bridge:deltaswapBridge,
    message:deltaswapMessage,
    emitter: deltaswapEmitter,
    sequence: deltaswapSequence,
    feeCollector: deltaswapFeeCollector,
    clock,
    rent,
    systemProgram,
  } = getPostMessageAccounts(
    deltaswapProgramId,
    payer,
    nftBridgeProgramId,
    message
  );
  return {
    payer: new PublicKey(payer),
    config: deriveNftBridgeConfigKey(nftBridgeProgramId),
    from: new PublicKey(from),
    mint: new PublicKey(mint),
    splMetadata: deriveTokenMetadataKey(mint),
    custody: deriveCustodyKey(nftBridgeProgramId, mint),
    authoritySigner: deriveAuthoritySignerKey(nftBridgeProgramId),
    custodySigner: deriveCustodySignerKey(nftBridgeProgramId),
   deltaswapBridge,
   deltaswapMessage,
    deltaswapEmitter,
    deltaswapSequence,
    deltaswapFeeCollector,
    clock,
    rent,
    systemProgram,
    tokenProgram: TOKEN_PROGRAM_ID,
    splMetadataProgram: TOKEN_METADATA_PROGRAM_ID,
    deltaswapProgram: new PublicKey(deltaswapProgramId),
  };
}
