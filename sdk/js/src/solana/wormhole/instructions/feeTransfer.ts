import {
  Commitment,
  Connection,
  PublicKey,
  PublicKeyInitData,
  SystemProgram,
  TransactionInstruction,
} from "@solana/web3.js";
import { deriveFeeCollectorKey, getDeltaswapBridgeData } from "../accounts";

export async function createBridgeFeeTransferInstruction(
  connection: Connection,
  deltaswapProgramId: PublicKeyInitData,
  payer: PublicKeyInitData,
  commitment?: Commitment
): Promise<TransactionInstruction> {
  const fee = await getDeltaswapBridgeData(
    connection,
    deltaswapProgramId,
    commitment
  ).then((data) => data.config.fee);
  return SystemProgram.transfer({
    fromPubkey: new PublicKey(payer),
    toPubkey: deriveFeeCollectorKey(deltaswapProgramId),
    lamports: fee,
  });
}
