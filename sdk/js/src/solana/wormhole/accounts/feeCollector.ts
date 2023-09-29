import { PublicKey, PublicKeyInitData } from "@solana/web3.js";
import { deriveAddress } from "../../utils";

export function deriveFeeCollectorKey(
  deltaswapProgramId: PublicKeyInitData
): PublicKey {
  return deriveAddress([Buffer.from("fee_collector")], deltaswapProgramId);
}
