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
  deriveCustodySignerKey,
  deriveTokenBridgeConfigKey,
  deriveCustodyKey,
} from "../accounts";

export function createTransferNativeInstruction(
  tokenBridgeProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  message: PublicKeyInitData,
  from: PublicKeyInitData,
  mint: PublicKeyInitData,
  nonce: number,
  amount: bigint,
  fee: bigint,
  targetAddress: Buffer | Uint8Array,
  targetChain: number
): TransactionInstruction {
  const methods = createReadOnlyTokenBridgeProgramInterface(
    tokenBridgeProgramId
  ).methods.transferNative(
    nonce,
    amount as any,
    fee as any,
    Buffer.from(targetAddress) as any,
    targetChain
  );

  // @ts-ignore
  return methods._ixFn(...methods._args, {
    accounts: getTransferNativeAccounts(
      tokenBridgeProgramId,
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
  deltaswapProgram: PublicKey;
}

export function getTransferNativeAccounts(
  tokenBridgeProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  message: PublicKeyInitData,
  from: PublicKeyInitData,
  mint: PublicKeyInitData
): TransferNativeAccounts {
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
    mint: new PublicKey(mint),
    custody: deriveCustodyKey(tokenBridgeProgramId, mint),
    authoritySigner: deriveAuthoritySignerKey(tokenBridgeProgramId),
    custodySigner: deriveCustodySignerKey(tokenBridgeProgramId),
   deltaswapBridge,
   deltaswapMessage:deltaswapMessage,
    deltaswapEmitter,
    deltaswapSequence,
    deltaswapFeeCollector,
    clock,
    rent,
    systemProgram,
    tokenProgram: TOKEN_PROGRAM_ID,
    deltaswapProgram: new PublicKey(deltaswapProgramId),
  };
}
