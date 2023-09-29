import {
  PublicKey,
  PublicKeyInitData,
  SystemProgram,
  SYSVAR_CLOCK_PUBKEY,
  SYSVAR_RENT_PUBKEY,
  TransactionInstruction,
} from "@solana/web3.js";
import {
  isBytes,
  ParsedGovernanceVaa,
  parseGovernanceVaa,
  SignedVaa,
} from "../../../vaa";
import { createReadOnlyDeltaswapProgramInterface } from "../program";
import {
  deriveDeltaswapBridgeDataKey,
  deriveClaimKey,
  deriveFeeCollectorKey,
  derivePhylaxSetKey,
  derivePostedVaaKey,
  deriveUpgradeAuthorityKey,
} from "../accounts";
import { BpfLoaderUpgradeable, deriveUpgradeableProgramKey } from "../../utils";

export function createSetFeesInstruction(
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  vaa: SignedVaa | ParsedGovernanceVaa
): TransactionInstruction {
  const methods =
    createReadOnlyDeltaswapProgramInterface(deltaswapProgramId).methods.setFees();

  // @ts-ignore
  return methods._ixFn(...methods._args, {
    accounts: getSetFeesAccounts(deltaswapProgramId, payer, vaa) as any,
    signers: undefined,
    remainingAccounts: undefined,
    preInstructions: undefined,
    postInstructions: undefined,
  });
}

export interface SetFeesAccounts {
  payer: PublicKey;
  bridge: PublicKey;
  vaa: PublicKey;
  claim: PublicKey;
  systemProgram: PublicKey;
}

export function getSetFeesAccounts(
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  vaa: SignedVaa | ParsedGovernanceVaa
): SetFeesAccounts {
  const parsed = isBytes(vaa) ? parseGovernanceVaa(vaa) : vaa;
  return {
    payer: new PublicKey(payer),
    bridge: deriveDeltaswapBridgeDataKey(deltaswapProgramId),
    vaa: derivePostedVaaKey(deltaswapProgramId, parsed.hash),
    claim: deriveClaimKey(
      deltaswapProgramId,
      parsed.emitterAddress,
      parsed.emitterChain,
      parsed.sequence
    ),
    systemProgram: SystemProgram.programId,
  };
}

export function createTransferFeesInstruction(
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  recipient: PublicKeyInitData,
  vaa: SignedVaa | ParsedGovernanceVaa
): TransactionInstruction {
  const methods =
    createReadOnlyDeltaswapProgramInterface(
      deltaswapProgramId
    ).methods.transferFees();

  // @ts-ignore
  return methods._ixFn(...methods._args, {
    accounts: getTransferFeesAccounts(
      deltaswapProgramId,
      payer,
      recipient,
      vaa
    ) as any,
    signers: undefined,
    remainingAccounts: undefined,
    preInstructions: undefined,
    postInstructions: undefined,
  });
}

export interface TransferFeesAccounts {
  payer: PublicKey;
  bridge: PublicKey;
  vaa: PublicKey;
  claim: PublicKey;
  feeCollector: PublicKey;
  recipient: PublicKey;
  rent: PublicKey;
  systemProgram: PublicKey;
}

export function getTransferFeesAccounts(
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  recipient: PublicKeyInitData,
  vaa: SignedVaa | ParsedGovernanceVaa
): TransferFeesAccounts {
  const parsed = isBytes(vaa) ? parseGovernanceVaa(vaa) : vaa;
  return {
    payer: new PublicKey(payer),
    bridge: deriveDeltaswapBridgeDataKey(deltaswapProgramId),
    vaa: derivePostedVaaKey(deltaswapProgramId, parsed.hash),
    claim: deriveClaimKey(
      deltaswapProgramId,
      parsed.emitterAddress,
      parsed.emitterChain,
      parsed.sequence
    ),
    feeCollector: deriveFeeCollectorKey(deltaswapProgramId),
    recipient: new PublicKey(recipient),
    rent: SYSVAR_RENT_PUBKEY,
    systemProgram: SystemProgram.programId,
  };
}

export function createUpgradePhylaxSetInstruction(
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  vaa: SignedVaa | ParsedGovernanceVaa
): TransactionInstruction {
  const methods =
    createReadOnlyDeltaswapProgramInterface(
      deltaswapProgramId
    ).methods.upgradePhylaxSet();

  // @ts-ignore
  return methods._ixFn(...methods._args, {
    accounts: getUpgradePhylaxSetAccounts(
      deltaswapProgramId,
      payer,
      vaa
    ) as any,
    signers: undefined,
    remainingAccounts: undefined,
    preInstructions: undefined,
    postInstructions: undefined,
  });
}

export interface UpgradePhylaxSetAccounts {
  payer: PublicKey;
  bridge: PublicKey;
  vaa: PublicKey;
  claim: PublicKey;
  phylaxSetOld: PublicKey;
  phylaxSetNew: PublicKey;
  systemProgram: PublicKey;
}

export function getUpgradePhylaxSetAccounts(
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  vaa: SignedVaa | ParsedGovernanceVaa
): UpgradePhylaxSetAccounts {
  const parsed = isBytes(vaa) ? parseGovernanceVaa(vaa) : vaa;
  return {
    payer: new PublicKey(payer),
    bridge: deriveDeltaswapBridgeDataKey(deltaswapProgramId),
    vaa: derivePostedVaaKey(deltaswapProgramId, parsed.hash),
    claim: deriveClaimKey(
      deltaswapProgramId,
      parsed.emitterAddress,
      parsed.emitterChain,
      parsed.sequence
    ),
    phylaxSetOld: derivePhylaxSetKey(
      deltaswapProgramId,
      parsed.phylaxSetIndex
    ),
    phylaxSetNew: derivePhylaxSetKey(
      deltaswapProgramId,
      parsed.phylaxSetIndex + 1
    ),
    systemProgram: SystemProgram.programId,
  };
}

export function createUpgradeContractInstruction(
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  vaa: SignedVaa | ParsedGovernanceVaa
): TransactionInstruction {
  const methods =
    createReadOnlyDeltaswapProgramInterface(
      deltaswapProgramId
    ).methods.upgradeContract();

  // @ts-ignore
  return methods._ixFn(...methods._args, {
    accounts: getUpgradeContractAccounts(deltaswapProgramId, payer, vaa) as any,
    signers: undefined,
    remainingAccounts: undefined,
    preInstructions: undefined,
    postInstructions: undefined,
  });
}

export interface UpgradeContractAccounts {
  payer: PublicKey;
  bridge: PublicKey;
  vaa: PublicKey;
  claim: PublicKey;
  upgradeAuthority: PublicKey;
  spill: PublicKey;
  implementation: PublicKey;
  programData: PublicKey;
  deltaswapProgram: PublicKey;
  rent: PublicKey;
  clock: PublicKey;
  bpfLoaderUpgradeable: PublicKey;
  systemProgram: PublicKey;
}

export function getUpgradeContractAccounts(
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  vaa: SignedVaa | ParsedGovernanceVaa,
  spill?: PublicKeyInitData
): UpgradeContractAccounts {
  const parsed = isBytes(vaa) ? parseGovernanceVaa(vaa) : vaa;
  const implementation = parsed.orderPayload;
  if (implementation.length != 32) {
    throw new Error("implementation.length != 32");
  }
  return {
    payer: new PublicKey(payer),
    bridge: deriveDeltaswapBridgeDataKey(deltaswapProgramId),
    vaa: derivePostedVaaKey(deltaswapProgramId, parsed.hash),
    claim: deriveClaimKey(
      deltaswapProgramId,
      parsed.emitterAddress,
      parsed.emitterChain,
      parsed.sequence
    ),
    upgradeAuthority: deriveUpgradeAuthorityKey(deltaswapProgramId),
    spill: new PublicKey(spill === undefined ? payer : spill),
    implementation: new PublicKey(implementation),
    programData: deriveUpgradeableProgramKey(deltaswapProgramId),
    deltaswapProgram: new PublicKey(deltaswapProgramId),
    rent: SYSVAR_RENT_PUBKEY,
    clock: SYSVAR_CLOCK_PUBKEY,
    bpfLoaderUpgradeable: BpfLoaderUpgradeable.programId,
    systemProgram: SystemProgram.programId,
  };
}
