import { PublicKeyInitData } from "@solana/web3.js";
import { deriveClaimKey } from "../solana/deltaswap";
import { parseVaa, SignedVaa } from "../vaa/deltaswap";

export async function getClaimAddressSolana(
  programAddress: PublicKeyInitData,
  signedVaa: SignedVaa
) {
  const parsed = parseVaa(signedVaa);
  return deriveClaimKey(
    programAddress,
    parsed.emitterAddress,
    parsed.emitterChain,
    parsed.sequence
  );
}
