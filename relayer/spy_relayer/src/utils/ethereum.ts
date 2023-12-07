import { TokenImplementation__factory } from "@deltaswapio/deltaswap-sdk/lib/cjs/ethers-contracts";
import { ethers } from "ethers";

export async function getEthereumToken(
  tokenAddress: string,
  provider: ethers.providers.Provider
) {
  const token = TokenImplementation__factory.connect(tokenAddress, provider);
  return token;
}
