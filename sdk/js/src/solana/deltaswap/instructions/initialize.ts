import {
  PublicKey,
  PublicKeyInitData,
  SystemProgram,
  SYSVAR_CLOCK_PUBKEY,
  SYSVAR_RENT_PUBKEY,
  TransactionInstruction,
} from "@solana/web3.js";
import { createReadOnlyDeltaswapProgramInterface } from "../program";
import {
  deriveFeeCollectorKey,
  derivePhylaxSetKey,
  deriveDeltaswapBridgeDataKey,
} from "../accounts";
import BN from "bn.js";

export function createInitializeInstruction(
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  phylaxSetExpirationTime: number,
  fee: bigint,
  initialPhylaxs: Buffer[]
): TransactionInstruction {
  const methods = createReadOnlyDeltaswapProgramInterface(
    deltaswapProgramId
  ).methods.initialize(phylaxSetExpirationTime, new BN(fee.toString()), [
    ...initialPhylaxs,
  ]);

  // @ts-ignore
  return methods._ixFn(...methods._args, {
    accounts: getInitializeAccounts(deltaswapProgramId, payer) as any,
    signers: undefined,
    remainingAccounts: undefined,
    preInstructions: undefined,
    postInstructions: undefined,
  });
}

export interface InitializeAccounts {
  bridge: PublicKey;
  phylaxSet: PublicKey;
  feeCollector: PublicKey;
  payer: PublicKey;
  clock: PublicKey;
  rent: PublicKey;
  systemProgram: PublicKey;
}

export function getInitializeAccounts(
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData
): InitializeAccounts {
  return {
    bridge: deriveDeltaswapBridgeDataKey(deltaswapProgramId),
    phylaxSet: derivePhylaxSetKey(deltaswapProgramId, 0),
    feeCollector: deriveFeeCollectorKey(deltaswapProgramId),
    payer: new PublicKey(payer),
    clock: SYSVAR_CLOCK_PUBKEY,
    rent: SYSVAR_RENT_PUBKEY,
    systemProgram: SystemProgram.programId,
  };
}
