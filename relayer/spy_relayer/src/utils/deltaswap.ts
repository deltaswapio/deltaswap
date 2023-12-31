import { ChainId } from "@deltaswapio/deltaswap-sdk";

export const chainIDStrings: { [key in ChainId]: string } = {
  0: "unset",
  1: "solana",
  2: "ethereum",
  3: "terra",
  4: "bsc",
  5: "polygon",
  6: "avalanche",
  7: "oasis",
  8: "algorand",
  9: "aurora",
  10: "fantom",
  11: "karura",
  12: "acala",
  13: "klaytn",
  14: "celo",
  15: "near",
  16: "moonbeam",
  17: "neon",
  18: "terra2",
  19: "injective",
  20: "osmosis",
  21: "sui",
  22: "aptos",
  23: "arbitrum",
  24: "optimism",
  25: "gnosis",
  26: "pythnet",
  28: "xpla",
  29: "btc",
  30: "base",
  32: "sei",
  33: "rootstock",
  34: "scroll",
  100: "tron",
  7070: "planq",
  7077: "deltachain",
  4000: "cosmoshub",
  4001: "evmos",
  4002: "kujira",
  4004: "celestia",
  10002: "sepolia",
};
