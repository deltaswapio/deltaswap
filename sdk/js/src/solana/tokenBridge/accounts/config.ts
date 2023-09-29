import {
  Connection,
  PublicKey,
  Commitment,
  PublicKeyInitData,
} from "@solana/web3.js";
import { deriveAddress, getAccountData } from "../../utils";

export function deriveTokenBridgeConfigKey(
  tokenBridgeProgramId: PublicKeyInitData
): PublicKey {
  return deriveAddress([Buffer.from("config")], tokenBridgeProgramId);
}

export async function getTokenBridgeConfig(
  connection: Connection,
  tokenBridgeProgramId: PublicKeyInitData,
  commitment?: Commitment
): Promise<TokenBridgeConfig> {
  return connection
    .getAccountInfo(
      deriveTokenBridgeConfigKey(tokenBridgeProgramId),
      commitment
    )
    .then((info) => TokenBridgeConfig.deserialize(getAccountData(info)));
}

export class TokenBridgeConfig {
  deltaswap: PublicKey;

  constructor(deltaswapProgramId: Buffer) {
    this.deltaswap = new PublicKey(deltaswapProgramId);
  }

  static deserialize(data: Buffer): TokenBridgeConfig {
    if (data.length != 32) {
      throw new Error("data.length != 32");
    }
    const deltaswapProgramId = data.subarray(0, 32);
    return new TokenBridgeConfig(deltaswapProgramId);
  }
}
