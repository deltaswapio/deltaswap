import {
  PublicKey,
  PublicKeyInitData,
  TransactionInstruction,
} from "@solana/web3.js";
import { TOKEN_PROGRAM_ID } from "@solana/spl-token";
import { createReadOnlyTokenBridgeProgramInterface } from "../program";
import { getPostMessageCpiAccounts } from "../../deltaswap";
import {
  deriveAuthoritySignerKey,
  deriveSenderAccountKey,
  deriveTokenBridgeConfigKey,
  deriveWrappedMetaKey,
  deriveWrappedMintKey,
} from "../accounts";

export function createTransferWrappedWithPayloadInstruction(
  tokenBridgeProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  message: PublicKeyInitData,
  from: PublicKeyInitData,
  fromOwner: PublicKeyInitData,
  tokenChain: number,
  tokenAddress: Buffer | Uint8Array,
  nonce: number,
  amount: bigint,
  targetAddress: Buffer | Uint8Array,
  targetChain: number,
  payload: Buffer | Uint8Array
): TransactionInstruction {
  const methods = createReadOnlyTokenBridgeProgramInterface(
    tokenBridgeProgramId
  ).methods.transferWrappedWithPayload(
    nonce,
    amount as any,
    Buffer.from(targetAddress) as any,
    targetChain,
    Buffer.from(payload) as any,
    null
  );

  // @ts-ignore
  return methods._ixFn(...methods._args, {
    accounts: getTransferWrappedWithPayloadAccounts(
      tokenBridgeProgramId,
      deltaswapProgramId,
      payer,
      message,
      from,
      fromOwner,
      tokenChain,
      tokenAddress
    ) as any,
    signers: undefined,
    remainingAccounts: undefined,
    preInstructions: undefined,
    postInstructions: undefined,
  });
}

export interface TransferWrappedWithPayloadAccounts {
  payer: PublicKey;
  config: PublicKey;
  from: PublicKey;
  fromOwner: PublicKey;
  mint: PublicKey;
  wrappedMeta: PublicKey;
  authoritySigner: PublicKey;
 deltaswapBridge: PublicKey;
 deltaswapMessage: PublicKey;
  deltaswapEmitter: PublicKey;
  deltaswapSequence: PublicKey;
  deltaswapFeeCollector: PublicKey;
  clock: PublicKey;
  sender: PublicKey;
  rent: PublicKey;
  systemProgram: PublicKey;
  tokenProgram: PublicKey;
  deltaswapProgram: PublicKey;
}

export function getTransferWrappedWithPayloadAccounts(
  tokenBridgeProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  message: PublicKeyInitData,
  from: PublicKeyInitData,
  fromOwner: PublicKeyInitData,
  tokenChain: number,
  tokenAddress: Buffer | Uint8Array,
  cpiProgramId?: PublicKeyInitData
): TransferWrappedWithPayloadAccounts {
  const mint = deriveWrappedMintKey(
    tokenBridgeProgramId,
    tokenChain,
    tokenAddress
  );
  const {
   deltaswapBridge,
   deltaswapMessage,
    deltaswapEmitter,
    deltaswapSequence,
    deltaswapFeeCollector,
    clock,
    rent,
    systemProgram,
  } = getPostMessageCpiAccounts(
    tokenBridgeProgramId,
    deltaswapProgramId,
    payer,
    message
  );
  return {
    payer: new PublicKey(payer),
    config: deriveTokenBridgeConfigKey(tokenBridgeProgramId),
    from: new PublicKey(from),
    fromOwner: new PublicKey(fromOwner),
    mint: mint,
    wrappedMeta: deriveWrappedMetaKey(tokenBridgeProgramId, mint),
    authoritySigner: deriveAuthoritySignerKey(tokenBridgeProgramId),
   deltaswapBridge,
   deltaswapMessage:deltaswapMessage,
    deltaswapEmitter,
    deltaswapSequence,
    deltaswapFeeCollector,
    clock,
    sender: new PublicKey(
      cpiProgramId === undefined ? payer : deriveSenderAccountKey(cpiProgramId)
    ),
    rent,
    systemProgram,
    deltaswapProgram: new PublicKey(deltaswapProgramId),
    tokenProgram: TOKEN_PROGRAM_ID,
  };
}
