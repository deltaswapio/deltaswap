import {
  Commitment,
  Connection,
  PublicKey,
  PublicKeyInitData,
} from "@solana/web3.js";
import { deriveAddress } from "../../utils";
import {
  deriveEmitterSequenceKey,
  getSequenceTracker,
  SequenceTracker,
} from "./sequence";

export interface EmitterAccounts {
  emitter: PublicKey;
  sequence: PublicKey;
}

export function deriveDeltaswapEmitterKey(
  emitterProgramId: PublicKeyInitData
): PublicKey {
  return deriveAddress([Buffer.from("emitter")], emitterProgramId);
}

export function getEmitterKeys(
  emitterProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData
): EmitterAccounts {
  const emitter = deriveDeltaswapEmitterKey(emitterProgramId);
  return {
    emitter,
    sequence: deriveEmitterSequenceKey(emitter, deltaswapProgramId),
  };
}

export async function getProgramSequenceTracker(
  connection: Connection,
  emitterProgramId: PublicKeyInitData,
  deltaswapProgramId: PublicKeyInitData,
  commitment?: Commitment
): Promise<SequenceTracker> {
  return getSequenceTracker(
    connection,
    deriveDeltaswapEmitterKey(emitterProgramId),
    deltaswapProgramId,
    commitment
  );
}
