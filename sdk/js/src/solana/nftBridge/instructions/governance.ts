import {
  PublicKey,
  PublicKeyInitData,
  SystemProgram,
  SYSVAR_CLOCK_PUBKEY,
  SYSVAR_RENT_PUBKEY,
  TransactionInstruction,
} from "@solana/web3.js";
import { createReadOnlyNftBridgeProgramInterface } from "../program";
import { deriveClaimKey, derivePostedVaaKey } from "../../deltaswap";
import {
  deriveEndpointKey,
  deriveNftBridgeConfigKey,
  deriveUpgradeAuthorityKey,
} from "../accounts";
import {
  isBytes,
  ParsedNftBridgeRegisterChainVaa,
  ParsedNftBridgeUpgradeContractVaa,
  parseNftBridgeRegisterChainVaa,
  parseNftBridgeUpgradeContractVaa,
  SignedVaa,
} from "../../../vaa";
import { BpfLoaderUpgradeable, deriveUpgradeableProgramKey } from "../../utils";

export function createRegisterChainInstruction(
  nftBridgeProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  vaa: SignedVaa | ParsedNftBridgeRegisterChainVaa
): TransactionInstruction {
  const methods =
    createReadOnlyNftBridgeProgramInterface(
      nftBridgeProgramId
    ).methods.registerChain();

  // @ts-ignore
  return methods._ixFn(...methods._args, {
    accounts: getRegisterChainAccounts(
      nftBridgeProgramId,
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

export interface RegisterChainAccounts {
  payer: PublicKey;
  config: PublicKey;
  endpoint: PublicKey;
  vaa: PublicKey;
  claim: PublicKey;
  rent: PublicKey;
  systemProgram: PublicKey;
  deltaswapProgram: PublicKey;
}

export function getRegisterChainAccounts(
  nftBridgeProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  vaa: SignedVaa | ParsedNftBridgeRegisterChainVaa
): RegisterChainAccounts {
  const parsed = isBytes(vaa) ? parseNftBridgeRegisterChainVaa(vaa) : vaa;
  return {
    payer: new PublicKey(payer),
    config: deriveNftBridgeConfigKey(nftBridgeProgramId),
    endpoint: deriveEndpointKey(
      nftBridgeProgramId,
      parsed.foreignChain,
      parsed.foreignAddress
    ),
    vaa: derivePostedVaaKey(deltaswapProgramId, parsed.hash),
    claim: deriveClaimKey(
      nftBridgeProgramId,
      parsed.emitterAddress,
      parsed.emitterChain,
      parsed.sequence
    ),
    rent: SYSVAR_RENT_PUBKEY,
    systemProgram: SystemProgram.programId,
    deltaswapProgram: new PublicKey(deltaswapProgramId),
  };
}

export function createUpgradeContractInstruction(
  nftBridgeProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  vaa: SignedVaa | ParsedNftBridgeUpgradeContractVaa,
  spill?: PublicKeyInitData
): TransactionInstruction {
  const methods =
    createReadOnlyNftBridgeProgramInterface(
      nftBridgeProgramId
    ).methods.upgradeContract();

  // @ts-ignore
  return methods._ixFn(...methods._args, {
    accounts: getUpgradeContractAccounts(
      nftBridgeProgramId,
      deltaswapProgramId,
      payer,
      vaa,
      spill
    ) as any,
    signers: undefined,
    remainingAccounts: undefined,
    preInstructions: undefined,
    postInstructions: undefined,
  });
}

export interface UpgradeContractAccounts {
  payer: PublicKey;
  vaa: PublicKey;
  claim: PublicKey;
  upgradeAuthority: PublicKey;
  spill: PublicKey;
  implementation: PublicKey;
  programData: PublicKey;
  nftBridgeProgram: PublicKey;
  rent: PublicKey;
  clock: PublicKey;
  bpfLoaderUpgradeable: PublicKey;
  systemProgram: PublicKey;
}

export function getUpgradeContractAccounts(
  nftBridgeProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  vaa: SignedVaa | ParsedNftBridgeUpgradeContractVaa,
  spill?: PublicKeyInitData
): UpgradeContractAccounts {
  const parsed = isBytes(vaa) ? parseNftBridgeUpgradeContractVaa(vaa) : vaa;
  return {
    payer: new PublicKey(payer),
    vaa: derivePostedVaaKey(deltaswapProgramId, parsed.hash),
    claim: deriveClaimKey(
      nftBridgeProgramId,
      parsed.emitterAddress,
      parsed.emitterChain,
      parsed.sequence
    ),
    upgradeAuthority: deriveUpgradeAuthorityKey(nftBridgeProgramId),
    spill: new PublicKey(spill === undefined ? payer : spill),
    implementation: new PublicKey(parsed.newContract),
    programData: deriveUpgradeableProgramKey(nftBridgeProgramId),
    nftBridgeProgram: new PublicKey(nftBridgeProgramId),
    rent: SYSVAR_RENT_PUBKEY,
    clock: SYSVAR_CLOCK_PUBKEY,
    bpfLoaderUpgradeable: BpfLoaderUpgradeable.programId,
    systemProgram: SystemProgram.programId,
  };
}
