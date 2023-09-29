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
  deriveNftBridgeConfigKey,
  deriveWrappedMetaKey,
  deriveWrappedMintKey,
} from "../accounts";
import { createReadOnlyNftBridgeProgramInterface } from "../program";

export function createTransferWrappedInstruction(
  nftBridgeProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  message: PublicKeyInitData,
  from: PublicKeyInitData,
  fromOwner: PublicKeyInitData,
  tokenChain: number,
  tokenAddress: Buffer | Uint8Array,
  tokenId: bigint | number,
  nonce: number,
  targetAddress: Buffer | Uint8Array,
  targetChain: number
): TransactionInstruction {
  const methods = createReadOnlyNftBridgeProgramInterface(
    nftBridgeProgramId
  ).methods.transferWrapped(
    nonce,
    Buffer.from(targetAddress) as any,
    targetChain
  );

  // @ts-ignore
  return methods._ixFn(...methods._args, {
    accounts: getTransferWrappedAccounts(
      nftBridgeProgramId,
      deltaswapProgramId,
      payer,
      message,
      from,
      fromOwner,
      tokenChain,
      tokenAddress,
      tokenId
    ) as any,
    signers: undefined,
    remainingAccounts: undefined,
    preInstructions: undefined,
    postInstructions: undefined,
  });
}

export interface TransferWrappedAccounts {
  payer: PublicKey;
  config: PublicKey;
  from: PublicKey;
  fromOwner: PublicKey;
  mint: PublicKey;
  wrappedMeta: PublicKey;
  splMetadata: PublicKey;
  authoritySigner: PublicKey;
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

export function getTransferWrappedAccounts(
  nftBridgeProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  message: PublicKeyInitData,
  from: PublicKeyInitData,
  fromOwner: PublicKeyInitData,
  tokenChain: number,
  tokenAddress: Buffer | Uint8Array,
  tokenId: bigint | number
): TransferWrappedAccounts {
  const mint = deriveWrappedMintKey(
    nftBridgeProgramId,
    tokenChain,
    tokenAddress,
    tokenId
  );
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
    fromOwner: new PublicKey(fromOwner),
    mint,
    wrappedMeta: deriveWrappedMetaKey(nftBridgeProgramId, mint),
    splMetadata: deriveTokenMetadataKey(mint),
    authoritySigner: deriveAuthoritySignerKey(nftBridgeProgramId),
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
