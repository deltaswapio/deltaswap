import {
  PublicKey,
  PublicKeyInitData,
  TransactionInstruction,
  SYSVAR_CLOCK_PUBKEY,
  SYSVAR_RENT_PUBKEY,
  SystemProgram,
} from "@solana/web3.js";
import { createReadOnlyDeltaswapProgramInterface } from "../program";
import {
  deriveDeltaswapBridgeDataKey,
  derivePhylaxSetKey,
  derivePostedVaaKey,
} from "../accounts";
import { isBytes, ParsedVaa, parseVaa, SignedVaa } from "../../../vaa";
import BN from "bn.js";

/**
 * Make {@link TransactionInstruction} for `post_vaa` instruction.
 *
 * This is used in {@link createPostSignedVaaTransactions}'s last transaction.
 * `signatureSet` is a {@link web3.Keypair} generated outside of this method, which was used
 * to write signatures and the message hash to.
 *
 * https://github.com/deltaswapio/deltaswap/blob/main/solana/bridge/program/src/api/post_vaa.rs
 *
 * @param {PublicKeyInitData} deltaswapProgramId - deltaswap program address
 * @param {PublicKeyInitData} payer - transaction signer address
 * @param {SignedVaa | ParsedVaa} vaa - either signed VAA bytes or parsed VAA (use {@link parseVaa} on signed VAA)
 * @param {PublicKeyInitData} signatureSet - key for signature set account
 */
export function createPostVaaInstruction(
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  vaa: SignedVaa | ParsedVaa,
  signatureSet: PublicKeyInitData
): TransactionInstruction {
  const parsed = isBytes(vaa) ? parseVaa(vaa) : vaa;
  const methods = createReadOnlyDeltaswapProgramInterface(
    deltaswapProgramId
  ).methods.postVaa(
    parsed.version,
    parsed.phylaxSetIndex,
    parsed.timestamp,
    parsed.nonce,
    parsed.emitterChain,
    [...parsed.emitterAddress],
    new BN(parsed.sequence.toString()),
    parsed.consistencyLevel,
    parsed.payload
  );

  // @ts-ignore
  return methods._ixFn(...methods._args, {
    accounts: getPostVaaAccounts(
      deltaswapProgramId,
      payer,
      signatureSet,
      parsed
    ) as any,
    signers: undefined,
    remainingAccounts: undefined,
    preInstructions: undefined,
    postInstructions: undefined,
  });
}

export interface PostVaaAccounts {
  phylaxSet: PublicKey;
  bridge: PublicKey;
  signatureSet: PublicKey;
  vaa: PublicKey;
  payer: PublicKey;
  clock: PublicKey;
  rent: PublicKey;
  systemProgram: PublicKey;
}

export function getPostVaaAccounts(
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  signatureSet: PublicKeyInitData,
  vaa: SignedVaa | ParsedVaa
): PostVaaAccounts {
  const parsed = isBytes(vaa) ? parseVaa(vaa) : vaa;
  return {
    phylaxSet: derivePhylaxSetKey(
      deltaswapProgramId,
      parsed.phylaxSetIndex
    ),
    bridge: deriveDeltaswapBridgeDataKey(deltaswapProgramId),
    signatureSet: new PublicKey(signatureSet),
    vaa: derivePostedVaaKey(deltaswapProgramId, parsed.hash),
    payer: new PublicKey(payer),
    clock: SYSVAR_CLOCK_PUBKEY,
    rent: SYSVAR_RENT_PUBKEY,
    systemProgram: SystemProgram.programId,
  };
}
