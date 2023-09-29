import {
  Commitment,
  Connection,
  PublicKey,
  PublicKeyInitData,
  SystemProgram,
  TransactionInstruction,
} from "@solana/web3.js";
import { Program, Provider } from "@project-serum/anchor";
import { createReadOnlyProvider } from "../utils";
import { DeltaswapCoder } from "./coder";
import { Deltaswap } from "../types/deltaswap";

import IDL from "../../anchor-idl/deltaswap.json";

export function createDeltaswapProgramInterface(
  programId: PublicKeyInitData,
  provider?: Provider
): Program<Deltaswap> {
  return new Program<Deltaswap>(
    IDL as Deltaswap,
    new PublicKey(programId),
    provider === undefined ? ({ connection: null } as any) : provider,
    coder()
  );
}

export function createReadOnlyDeltaswapProgramInterface(
  programId: PublicKeyInitData,
  connection?: Connection
): Program<Deltaswap> {
  return createDeltaswapProgramInterface(
    programId,
    createReadOnlyProvider(connection)
  );
}

export function coder(): DeltaswapCoder {
  return new DeltaswapCoder(IDL as Deltaswap);
}
