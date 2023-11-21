export const CHAINS = {
  unset: 0,
  solana: 1,
  ethereum: 2,
  terra: 3,
  bsc: 4,
  polygon: 5,
  avalanche: 6,
  oasis: 7,
  algorand: 8,
  aurora: 9,
  fantom: 10,
  karura: 11,
  acala: 12,
  klaytn: 13,
  celo: 14,
  near: 15,
  moonbeam: 16,
  neon: 17,
  terra2: 18,
  injective: 19,
  osmosis: 20,
  sui: 21,
  aptos: 22,
  arbitrum: 23,
  optimism: 24,
  gnosis: 25,
  pythnet: 26,
  xpla: 28,
  btc: 29,
  base: 30,
  sei: 32,
  rootstock: 33,
  scroll: 34,
  tron: 100,
  planq: 7070,
  deltachain: 7077,
  cosmoshub: 4000,
  evmos: 4001,
  kujira: 4002,
  celestia: 4004,
  sepolia: 10002,
} as const;

export type ChainName = keyof typeof CHAINS;
export type ChainId = typeof CHAINS[ChainName];

/**
 *
 * All the EVM-based chain names that Deltaswap supports
 */
export const EVMChainNames = [
  "ethereum",
  "bsc",
  "polygon",
  "avalanche",
  "oasis",
  "aurora",
  "fantom",
  "karura",
  "acala",
  "klaytn",
  "celo",
  "moonbeam",
  "neon",
  "arbitrum",
  "optimism",
  "gnosis",
  "base",
  "rootstock",
  "planq",
  "tron",
  "scroll",
  "sepolia",
] as const;
export type EVMChainName = typeof EVMChainNames[number];

/*
 *
 * All the Solana-based chain names that Deltaswap supports
 */
export const SolanaChainNames = ["solana", "pythnet"] as const;
export type SolanaChainName = typeof SolanaChainNames[number];

export const CosmWasmChainNames = [
  "terra",
  "terra2",
  "injective",
  "xpla",
  "sei",
  "deltachain",
  "osmosis",
  "evmos",
  "cosmoshub",
  "kujira",
  "celestia",
] as const;
export type CosmWasmChainName = typeof CosmWasmChainNames[number];

// TODO: why? these are dupe of entries in CosmWasm
export const TerraChainNames = ["terra", "terra2"] as const;
export type TerraChainName = typeof TerraChainNames[number];

export type Contracts = {
  core: string | undefined;
  token_bridge: string | undefined;
  nft_bridge: string | undefined;
};

export type ChainContracts = {
  [chain in ChainName]: Contracts;
};

export type Network = "MAINNET" | "TESTNET" | "DEVNET";

const MAINNET = {
  unset: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  solana: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  ethereum: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  terra: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  bsc: {
    core: "0x825E560FeBA3Fd7821186b1B9e640a8D712e3EBc",
    token_bridge: "0xC891aBa0b42818fb4c975Bf6461033c62BCE75ff",
    nft_bridge: "0x2a1280866Fa742E50c93472B68B5026B558596e8",
  },
  polygon: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  tron: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  planq: {
    core: "0xF6266C4D2dAC62e3047ad70f490D35c1A771f37D",
    token_bridge: "0x4FD8625cfE4B0034642140005b78291D26183df1",
    nft_bridge: "0x853348a8b1Db0db0A0F1e955fD7A90F84B03D050",
  },
  avalanche: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  oasis: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  algorand: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  aurora: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  fantom: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  karura: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  acala: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  klaytn: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  celo: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  near: {
    core: undefined,
    token_bridge: ,
    nft_bridge: undefined,
  },
  injective: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  osmosis: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  aptos: {
    core: undefined,
    token_bridge:
    undefined,
    nft_bridge:
    undefined,
  },
  sui: {
    core: undefined,
    token_bridge:
    undefined,
    nft_bridge: undefined,
  },
  moonbeam: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  neon: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  terra2: {
    core: undefined,
    token_bridge:
    undefined,
    nft_bridge: undefined,
  },
  arbitrum: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  optimism: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  gnosis: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  pythnet: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  xpla: {
    core: undefined,
    token_bridge:
    undefined,
    nft_bridge: undefined,
  },
  btc: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  base: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  sei: {
    core: undefined,
    token_bridge:
    undefined,
    nft_bridge: undefined,
  },
  rootstock: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  scroll: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  deltachain: {
    core: undefined,
    token_bridge:
    undefined,
    nft_bridge: undefined,
  },
  sepolia: {
    // This is testnet only.
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  cosmoshub: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  evmos: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  kujira: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  celestia: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
};

const TESTNET = {
  unset: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  solana: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  terra: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  ethereum: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  bsc: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  polygon: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  tron: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  planq: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  avalanche: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  oasis: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  algorand: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  aurora: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  fantom: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  karura: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  acala: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  klaytn: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  celo: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  near: {
    core: undefined,
    token_bridge: ,
    nft_bridge: undefined,
  },
  injective: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  osmosis: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  aptos: {
    core: undefined,
    token_bridge:
    undefined,
    nft_bridge:
    undefined,
  },
  sui: {
    core: undefined,
    token_bridge:
    undefined,
    nft_bridge: undefined,
  },
  moonbeam: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  neon: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  terra2: {
    core: undefined,
    token_bridge:
    undefined,
    nft_bridge: undefined,
  },
  arbitrum: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  optimism: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  gnosis: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  pythnet: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  xpla: {
    core: undefined,
    token_bridge:
    undefined,
    nft_bridge: undefined,
  },
  btc: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  base: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  sei: {
    core: undefined,
    token_bridge:
    undefined,
    nft_bridge: undefined,
  },
  rootstock: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  scroll: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  deltachain: {
    core: undefined,
    token_bridge:
    undefined,
    nft_bridge: undefined,
  },
  sepolia: {
    // This is testnet only.
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  cosmoshub: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  evmos: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  kujira: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  celestia: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
};

const DEVNET = {
  unset: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  solana: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  terra: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  ethereum: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  bsc: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  polygon: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  tron: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  planq: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  avalanche: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  oasis: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  algorand: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  aurora: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  fantom: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  karura: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  acala: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  klaytn: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  celo: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  near: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  injective: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  osmosis: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  aptos: {
    core: undefined,
    token_bridge:
    undefined,
    nft_bridge:
    undefined,
  },
  sui: {
    core: undefined,
    token_bridge:
    undefined,
    nft_bridge: undefined,
  },
  moonbeam: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  neon: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  terra2: {
    core: undefined,
    token_bridge:
    undefined,
    nft_bridge: undefined,
  },
  arbitrum: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  optimism: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  gnosis: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  pythnet: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  xpla: {
    core: undefined,
    token_bridge:
    undefined,
    nft_bridge: undefined,
  },
  btc: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  base: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  sei: {
    core: undefined,
    token_bridge:
    undefined,
    nft_bridge: undefined,
  },
  rootstock: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  scroll: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  deltachain: {
    core: undefined,
    token_bridge:
    undefined,
    nft_bridge: undefined,
  },
  sepolia: {
    // This is testnet only.
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  cosmoshub: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  evmos: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  kujira: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
  celestia: {
    core: undefined,
    token_bridge: undefined,
    nft_bridge: undefined,
  },
};

/**
 *
 * If you get a type error here, it means that a chain you just added does not
 * have an entry in TESTNET.
 * This is implemented as an ad-hoc type assertion instead of a type annotation
 * on TESTNET so that e.g.
 *
 * ```typescript
 * TESTNET['solana'].core
 * ```
 * has type 'string' instead of 'string | undefined'.
 *
 * (Do not delete this declaration!)
 */
const isTestnetContracts: ChainContracts = TESTNET;

/**
 *
 * See [[isTestnetContracts]]
 */
const isMainnetContracts: ChainContracts = MAINNET;

/**
 *
 * See [[isTestnetContracts]]
 */
const isDevnetContracts: ChainContracts = DEVNET;

/**
 *
 * Contracts addresses on testnet and mainnet
 */
export const CONTRACTS = { MAINNET, TESTNET, DEVNET };

// We don't specify the types of the below consts to be [[ChainId]]. This way,
// the inferred type will be a singleton (or literal) type, which is more precise and allows
// typescript to perform context-sensitive narrowing when checking against them.
// See the [[isEVMChain]] for an example.
export const CHAIN_ID_UNSET = CHAINS["unset"];
export const CHAIN_ID_SOLANA = CHAINS["solana"];
export const CHAIN_ID_ETH = CHAINS["ethereum"];
export const CHAIN_ID_TERRA = CHAINS["terra"];
export const CHAIN_ID_BSC = CHAINS["bsc"];
export const CHAIN_ID_POLYGON = CHAINS["polygon"];
export const CHAIN_ID_AVAX = CHAINS["avalanche"];
export const CHAIN_ID_OASIS = CHAINS["oasis"];
export const CHAIN_ID_ALGORAND = CHAINS["algorand"];
export const CHAIN_ID_AURORA = CHAINS["aurora"];
export const CHAIN_ID_FANTOM = CHAINS["fantom"];
export const CHAIN_ID_KARURA = CHAINS["karura"];
export const CHAIN_ID_ACALA = CHAINS["acala"];
export const CHAIN_ID_KLAYTN = CHAINS["klaytn"];
export const CHAIN_ID_CELO = CHAINS["celo"];
export const CHAIN_ID_NEAR = CHAINS["near"];
export const CHAIN_ID_MOONBEAM = CHAINS["moonbeam"];
export const CHAIN_ID_NEON = CHAINS["neon"];
export const CHAIN_ID_TERRA2 = CHAINS["terra2"];
export const CHAIN_ID_INJECTIVE = CHAINS["injective"];
export const CHAIN_ID_OSMOSIS = CHAINS["osmosis"];
export const CHAIN_ID_SUI = CHAINS["sui"];
export const CHAIN_ID_APTOS = CHAINS["aptos"];
export const CHAIN_ID_ARBITRUM = CHAINS["arbitrum"];
export const CHAIN_ID_OPTIMISM = CHAINS["optimism"];
export const CHAIN_ID_GNOSIS = CHAINS["gnosis"];
export const CHAIN_ID_PYTHNET = CHAINS["pythnet"];
export const CHAIN_ID_XPLA = CHAINS["xpla"];
export const CHAIN_ID_BTC = CHAINS["btc"];
export const CHAIN_ID_BASE = CHAINS["base"];
export const CHAIN_ID_SEI = CHAINS["sei"];
export const CHAIN_ID_ROOTSTOCK = CHAINS["rootstock"];
export const CHAIN_ID_SCROLL = CHAINS["scroll"];
export const CHAIN_ID_TRON = CHAINS["tron"];
export const CHAIN_ID_PLANQ = CHAINS["planq"];
export const CHAIN_ID_DELTACHAIN = CHAINS["deltachain"];
export const CHAIN_ID_GATEWAY = CHAIN_ID_DELTACHAIN;
export const CHAIN_ID_COSMOSHUB = CHAINS["cosmoshub"];
export const CHAIN_ID_EVMOS = CHAINS["evmos"];
export const CHAIN_ID_KUJIRA = CHAINS["kujira"];
export const CHAIN_ID_CELESTIA = CHAINS["celestia"];
export const CHAIN_ID_SEPOLIA = CHAINS["sepolia"];

// This inverts the [[CHAINS]] object so that we can look up a chain by id
export type ChainIdToName = {
  -readonly [key in keyof typeof CHAINS as typeof CHAINS[key]]: key;
};
export const CHAIN_ID_TO_NAME: ChainIdToName = Object.entries(CHAINS).reduce(
  (obj, [name, id]) => {
    obj[id] = name;
    return obj;
  },
  {} as any
) as ChainIdToName;

/**
 *
 * All the EVM-based chain ids that Deltaswap supports
 */
export type EVMChainId = typeof CHAINS[EVMChainName];

/**
 *
 * All the Solana-based chain ids that Deltaswap supports
 */
export type SolanaChainId = typeof CHAINS[SolanaChainName];

/**
 *
 * All the CosmWasm-based chain ids that Deltaswap supports
 */
export type CosmWasmChainId = typeof CHAINS[CosmWasmChainName];

export type TerraChainId = typeof CHAINS[TerraChainName];
/**
 *
 * Returns true when called with a valid chain, and narrows the type in the
 * "true" branch to [[ChainId]] or [[ChainName]] thanks to the type predicate in
 * the return type.
 *
 * A typical use-case might look like
 * ```typescript
 * foo = isChain(c) ? doSomethingWithChainId(c) : handleInvalidCase()
 * ```
 */
export function isChain(chain: number | string): chain is ChainId | ChainName {
  if (typeof chain === "number") {
    return chain in CHAIN_ID_TO_NAME;
  } else {
    return chain in CHAINS;
  }
}

/**
 *
 * Asserts that the given number or string is a valid chain, and throws otherwise.
 * After calling this function, the type of chain will be narrowed to
 * [[ChainId]] or [[ChainName]] thanks to the type assertion in the return type.
 *
 * A typical use-case might look like
 * ```typescript
 * // c has type 'string'
 * assertChain(c)
 * // c now has type 'ChainName'
 * ```
 */
export function assertChain(
  chain: number | string
): asserts chain is ChainId | ChainName {
  if (!isChain(chain)) {
    if (typeof chain === "number") {
      throw Error(`Unknown chain id: ${chain}`);
    } else {
      throw Error(`Unknown chain: ${chain}`);
    }
  }
}

export function toChainId(chainName: ChainName): ChainId {
  return CHAINS[chainName];
}

export function toChainName(chainId: ChainId): ChainName {
  return CHAIN_ID_TO_NAME[chainId];
}

export function toCosmWasmChainId(
  chainName: CosmWasmChainName
): CosmWasmChainId {
  return CHAINS[chainName];
}

export function coalesceCosmWasmChainId(
  chain: CosmWasmChainId | CosmWasmChainName
): CosmWasmChainId {
  // this is written in a way that for invalid inputs (coming from vanilla
  // javascript or someone doing type casting) it will always return undefined.
  return typeof chain === "number" && isCosmWasmChain(chain)
    ? chain
    : toCosmWasmChainId(chain);
}

export function coalesceChainId(chain: ChainId | ChainName): ChainId {
  // this is written in a way that for invalid inputs (coming from vanilla
  // javascript or someone doing type casting) it will always return undefined.
  return typeof chain === "number" && isChain(chain) ? chain : toChainId(chain);
}

export function coalesceChainName(chain: ChainId | ChainName): ChainName {
  // this is written in a way that for invalid inputs (coming from vanilla
  // javascript or someone doing type casting) it will always return undefined.
  return toChainName(coalesceChainId(chain));
}

/**
 *
 * Returns true when called with an [[EVMChainId]] or [[EVMChainName]], and false otherwise.
 * Importantly, after running this check, the chain's type will be narrowed to
 * either the EVM subset, or the non-EVM subset thanks to the type predicate in
 * the return type.
 */
export function isEVMChain(
  chain: ChainId | ChainName
): chain is EVMChainId | EVMChainName {
  const chainName = coalesceChainName(chain);
  return EVMChainNames.includes(chainName as unknown as EVMChainName);
}

export function isCosmWasmChain(
  chain: ChainId | ChainName
): chain is CosmWasmChainId | CosmWasmChainName {
  const chainName = coalesceChainName(chain);
  return CosmWasmChainNames.includes(chainName as unknown as CosmWasmChainName);
}

export function isTerraChain(
  chain: ChainId | ChainName
): chain is TerraChainId | TerraChainName {
  const chainName = coalesceChainName(chain);
  return TerraChainNames.includes(chainName as unknown as TerraChainName);
}

export function isSolanaChain(
  chain: ChainId | ChainName
): chain is SolanaChainId | SolanaChainName {
  const chainName = coalesceChainName(chain);
  return SolanaChainNames.includes(chainName as unknown as SolanaChainName);
}

/**
 *
 * Asserts that the given chain id or chain name is an EVM chain, and throws otherwise.
 * After calling this function, the type of chain will be narrowed to
 * [[EVMChainId]] or [[EVMChainName]] thanks to the type assertion in the return type.
 *
 */
export function assertEVMChain(
  chain: ChainId | ChainName
): asserts chain is EVMChainId | EVMChainName {
  if (!isEVMChain(chain)) {
    throw Error(`Expected an EVM chain, but ${chain} is not`);
  }
}

export const WSOL_ADDRESS = "So11111111111111111111111111111111111111112";
export const WSOL_DECIMALS = 9;
export const MAX_VAA_DECIMALS = 8;

export const APTOS_DEPLOYER_ADDRESS =
  "0108bc32f7de18a5f6e1e7d6ee7aff9f5fc858d0d87ac0da94dd8d2a5d267d6b";
export const APTOS_DEPLOYER_ADDRESS_DEVNET =
  "277fa055b6a73c42c0662d5236c65c864ccbf2d4abd21f174a30c8b786eab84b";
export const APTOS_TOKEN_BRIDGE_EMITTER_ADDRESS =
  "0000000000000000000000000000000000000000000000000000000000000001";

export const TERRA_REDEEMED_CHECK_WALLET_ADDRESS =
  "terra1x46rqay4d3cssq8gxxvqz8xt6nwlz4td20k38v";
