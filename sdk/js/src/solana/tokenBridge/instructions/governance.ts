import {
  PublicKey,
  PublicKeyInitData,
  SystemProgram,
  SYSVAR_CLOCK_PUBKEY,
  SYSVAR_RENT_PUBKEY,
  TransactionInstruction,
} from "@solana/web3.js";
import { createReadOnlyTokenBridgeProgramInterface } from "../program";
import { deriveClaimKey, derivePostedVaaKey } from "../../deltaswap";
import {
  deriveEndpointKey,
  deriveTokenBridgeConfigKey,
  deriveUpgradeAuthorityKey,
} from "../accounts";
import {
  isBytes,
  ParsedTokenBridgeRegisterChainVaa,
  ParsedTokenBridgeUpgradeContractVaa,
  parseTokenBridgeRegisterChainVaa,
  parseTokenBridgeUpgradeContractVaa,
  SignedVaa,
} from "../../../vaa";
import { BpfLoaderUpgradeable, deriveUpgradeableProgramKey } from "../../utils";

export function createRegisterChainInstruction(
  tokenBridgeProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  vaa: SignedVaa | ParsedTokenBridgeRegisterChainVaa
): TransactionInstruction {
  const methods =
    createReadOnlyTokenBridgeProgramInterface(
      tokenBridgeProgramId
    ).methods.registerChain();

  // @ts-ignore
  return methods._ixFn(...methods._args, {
    accounts: getRegisterChainAccounts(
      tokenBridgeProgramId,
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
  tokenBridgeProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  vaa: SignedVaa | ParsedTokenBridgeRegisterChainVaa
): RegisterChainAccounts {
  const parsed = isBytes(vaa) ? parseTokenBridgeRegisterChainVaa(vaa) : vaa;
  return {
    payer: new PublicKey(payer),
    config: deriveTokenBridgeConfigKey(tokenBridgeProgramId),
    endpoint: deriveEndpointKey(
      tokenBridgeProgramId,
      parsed.foreignChain,
      parsed.foreignAddress
    ),
    vaa: derivePostedVaaKey(deltaswapProgramId, parsed.hash),
    claim: deriveClaimKey(
      tokenBridgeProgramId,
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
  tokenBridgeProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  vaa: SignedVaa | ParsedTokenBridgeUpgradeContractVaa,
  spill?: PublicKeyInitData
): TransactionInstruction {
  const methods =
    createReadOnlyTokenBridgeProgramInterface(
      tokenBridgeProgramId
    ).methods.upgradeContract();

  // @ts-ignore
  return methods._ixFn(...methods._args, {
    accounts: getUpgradeContractAccounts(
      tokenBridgeProgramId,
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
  tokenBridgeProgram: PublicKey;
  rent: PublicKey;
  clock: PublicKey;
  bpfLoaderUpgradeable: PublicKey;
  systemProgram: PublicKey;
}

export function getUpgradeContractAccounts(
  tokenBridgeProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  vaa: SignedVaa | ParsedTokenBridgeUpgradeContractVaa,
  spill?: PublicKeyInitData
): UpgradeContractAccounts {
  const parsed = isBytes(vaa) ? parseTokenBridgeUpgradeContractVaa(vaa) : vaa;
  return {
    payer: new PublicKey(payer),
    vaa: derivePostedVaaKey(deltaswapProgramId, parsed.hash),
    claim: deriveClaimKey(
      tokenBridgeProgramId,
      parsed.emitterAddress,
      parsed.emitterChain,
      parsed.sequence
    ),
    upgradeAuthority: deriveUpgradeAuthorityKey(tokenBridgeProgramId),
    spill: new PublicKey(spill === undefined ? payer : spill),
    implementation: new PublicKey(parsed.newContract),
    programData: deriveUpgradeableProgramKey(tokenBridgeProgramId),
    tokenBridgeProgram: new PublicKey(tokenBridgeProgramId),
    rent: SYSVAR_RENT_PUBKEY,
    clock: SYSVAR_CLOCK_PUBKEY,
    bpfLoaderUpgradeable: BpfLoaderUpgradeable.programId,
    systemProgram: SystemProgram.programId,
  };
}
