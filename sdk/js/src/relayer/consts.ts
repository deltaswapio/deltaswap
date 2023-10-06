import { ChainId, Network, ChainName, CHAIN_ID_TO_NAME } from "../";
import { ethers } from "ethers";
import {
  DeltaswapRelayer__factory,
  DeltaswapRelayer,
} from "../ethers-contracts/";

type AddressInfo = {
  deltaswapRelayerAddress?: string;
  mockDeliveryProviderAddress?: string;
  mockIntegrationAddress?: string;
};

const TESTNET: { [K in ChainName]?: AddressInfo } = {
  bsc: {
    deltaswapRelayerAddress: "0x80aC94316391752A193C1c47E27D382b507c93F3",
    mockDeliveryProviderAddress: "0x60a86b97a7596eBFd25fb769053894ed0D9A8366",
    mockIntegrationAddress: "0xb6A04D6672F005787147472Be20d39741929Aa03",
  },
  polygon: {
    deltaswapRelayerAddress: "0x0591C25ebd0580E0d4F27A82Fc2e24E7489CB5e0",
    mockDeliveryProviderAddress: "0x60a86b97a7596eBFd25fb769053894ed0D9A8366",
    mockIntegrationAddress: "0x3bF0c43d88541BBCF92bE508ec41e540FbF28C56",
  },
  avalanche: {
    deltaswapRelayerAddress: "0xA3cF45939bD6260bcFe3D66bc73d60f19e49a8BB",
    mockDeliveryProviderAddress: "0x60a86b97a7596eBFd25fb769053894ed0D9A8366",
    mockIntegrationAddress: "0x5E52f3eB0774E5e5f37760BD3Fca64951D8F74Ae",
  },
  celo: {
    deltaswapRelayerAddress: "0x306B68267Deb7c5DfCDa3619E22E9Ca39C374f84",
    mockDeliveryProviderAddress: "0x60a86b97a7596eBFd25fb769053894ed0D9A8366",
    mockIntegrationAddress: "0x7f1d8E809aBB3F6Dc9B90F0131C3E8308046E190",
  },
  moonbeam: {
    deltaswapRelayerAddress: "0x0591C25ebd0580E0d4F27A82Fc2e24E7489CB5e0",
    mockDeliveryProviderAddress: "0x60a86b97a7596eBFd25fb769053894ed0D9A8366",
    mockIntegrationAddress: "0x3bF0c43d88541BBCF92bE508ec41e540FbF28C56",
  },
  base: {
    deltaswapRelayerAddress: "0xea8029CD7FCAEFFcD1F53686430Db0Fc8ed384E1",
    mockDeliveryProviderAddress: "0x60a86b97a7596eBFd25fb769053894ed0D9A8366",
    mockIntegrationAddress: "0x9Ee656203B0DC40cc1bA3f4738527779220e3998",
  },
};

const DEVNET: { [K in ChainName]?: AddressInfo } = {
  ethereum: {
    deltaswapRelayerAddress: "0xE66C1Bc1b369EF4F376b84373E3Aa004E8F4C083",
    mockDeliveryProviderAddress: "0x1ef9e15c3bbf0555860b5009B51722027134d53a",
    mockIntegrationAddress: "0x0eb0dD3aa41bD15C706BC09bC03C002b7B85aeAC",
  },
  bsc: {
    deltaswapRelayerAddress: "0xE66C1Bc1b369EF4F376b84373E3Aa004E8F4C083",
    mockDeliveryProviderAddress: "0x1ef9e15c3bbf0555860b5009B51722027134d53a",
    mockIntegrationAddress: "0x0eb0dD3aa41bD15C706BC09bC03C002b7B85aeAC",
  },
};

const MAINNET: { [K in ChainName]?: AddressInfo } = {
  ethereum: {
    deltaswapRelayerAddress: "0x27428DD2d3DD32A4D7f7C497eAaa23130d894911",
  },
  bsc: {
    deltaswapRelayerAddress: "0x27428DD2d3DD32A4D7f7C497eAaa23130d894911",
  },
  polygon: {
    deltaswapRelayerAddress: "0x27428DD2d3DD32A4D7f7C497eAaa23130d894911",
  },
  planq: {
    deltaswapRelayerAddress: "0x27428DD2d3DD32A4D7f7C497eAaa23130d894911",
  },
  avalanche: {
    deltaswapRelayerAddress: "0x27428DD2d3DD32A4D7f7C497eAaa23130d894911",
  },
  fantom: {
    deltaswapRelayerAddress: "0x27428DD2d3DD32A4D7f7C497eAaa23130d894911",
  },
  klaytn: {
    deltaswapRelayerAddress: "0x27428DD2d3DD32A4D7f7C497eAaa23130d894911",
  },
  celo: {
    deltaswapRelayerAddress: "0x27428DD2d3DD32A4D7f7C497eAaa23130d894911",
  },
  acala: {
    deltaswapRelayerAddress: "0x27428DD2d3DD32A4D7f7C497eAaa23130d894911",
  },
  karura: {
    deltaswapRelayerAddress: "0x27428DD2d3DD32A4D7f7C497eAaa23130d894911",
  },
  moonbeam: {
    deltaswapRelayerAddress: "0x27428DD2d3DD32A4D7f7C497eAaa23130d894911",
  },
  arbitrum: {
    deltaswapRelayerAddress: "0x27428DD2d3DD32A4D7f7C497eAaa23130d894911",
  },
  optimism: {
    deltaswapRelayerAddress: "0x27428DD2d3DD32A4D7f7C497eAaa23130d894911",
  },
  base: {
    deltaswapRelayerAddress: "0x706f82e9bb5b0813501714ab5974216704980e31",
  },
};

export const RELAYER_CONTRACTS = { MAINNET, TESTNET, DEVNET };

export function getAddressInfo(
  chainName: ChainName,
  env: Network
): AddressInfo {
  const result: AddressInfo | undefined = RELAYER_CONTRACTS[env][chainName];
  if (!result) throw Error(`No address info for chain ${chainName} on ${env}`);
  return result;
}

export function getDeltaswapRelayerAddress(
  chainName: ChainName,
  env: Network
): string {
  const result = getAddressInfo(chainName, env).deltaswapRelayerAddress;
  if (!result)
    throw Error(
      `No Deltaswap Relayer Address for chain ${chainName}, network ${env}`
    );
  return result;
}

export function getDeltaswapRelayer(
  chainName: ChainName,
  env: Network,
  provider: ethers.providers.Provider | ethers.Signer,
  deltaswapRelayerAddress?: string
): DeltaswapRelayer {
  const thisChainsRelayer =
    deltaswapRelayerAddress || getDeltaswapRelayerAddress(chainName, env);
  const contract = DeltaswapRelayer__factory.connect(
    thisChainsRelayer,
    provider
  );
  return contract;
}

export const RPCS_BY_CHAIN: {
  [key in Network]: { [key in ChainName]?: string };
} = {
  MAINNET: {
    ethereum: "https://rpc.ankr.com/eth",
    bsc: "https://bsc-dataseed2.defibit.io",
    polygon: "https://rpc.ankr.com/polygon",
    planq: "https://evm-rpc.planq.network",
    avalanche: "https://rpc.ankr.com/avalanche",
    oasis: "https://emerald.oasis.dev",
    algorand: "https://mainnet-api.algonode.cloud",
    fantom: "https://rpc.ankr.com/fantom",
    karura: "https://eth-rpc-karura.aca-api.network",
    acala: "https://eth-rpc-acala.aca-api.network",
    klaytn: "https://klaytn-mainnet-rpc.allthatnode.com:8551",
    celo: "https://forno.celo.org",
    moonbeam: "https://rpc.ankr.com/moonbeam",
    arbitrum: "https://rpc.ankr.com/arbitrum",
    optimism: "https://rpc.ankr.com/optimism",
    aptos: "https://fullnode.mainnet.aptoslabs.com/",
    near: "https://rpc.mainnet.near.org",
    xpla: "https://dimension-lcd.xpla.dev",
    sui: "https://fullnode.mainnet.sui.io:443",
    terra2: "https://phoenix-lcd.terra.dev",
    terra: "https://columbus-fcd.terra.dev",
    injective: "https://k8s.mainnet.lcd.injective.network",
    solana: "https://api.mainnet-beta.solana.com",
  },
  TESTNET: {
    solana: "https://api.devnet.solana.com",
    terra: "https://bombay-lcd.terra.dev",
    ethereum: "https://rpc.ankr.com/eth_goerli",
    bsc: "https://bsc-testnet.publicnode.com",
    polygon: "https://rpc.ankr.com/polygon_mumbai",
    avalanche: "https://rpc.ankr.com/avalanche_fuji",
    oasis: "https://testnet.emerald.oasis.dev",
    algorand: "https://testnet-api.algonode.cloud",
    fantom: "https://rpc.testnet.fantom.network",
    aurora: "https://testnet.aurora.dev",
    karura: "https://karura-dev.aca-dev.network/eth/http",
    acala: "https://acala-dev.aca-dev.network/eth/http",
    klaytn: "https://api.baobab.klaytn.net:8651",
    celo: "https://alfajores-forno.celo-testnet.org",
    near: "https://rpc.testnet.near.org",
    injective: "https://k8s.testnet.tm.injective.network:443",
    aptos: "https://fullnode.testnet.aptoslabs.com/v1",
    pythnet: "https://api.pythtest.pyth.network/",
    xpla: "https://cube-lcd.xpla.dev:443",
    moonbeam: "https://rpc.api.moonbase.moonbeam.network",
    neon: "https://proxy.devnet.neonlabs.org/solana",
    terra2: "https://pisco-lcd.terra.dev",
    arbitrum: "https://goerli-rollup.arbitrum.io/rpc",
    optimism: "https://goerli.optimism.io",
    gnosis: "https://sokol.poa.network/",
    rootstock: "https://public-node.rsk.co",
  },
  DEVNET: {
    ethereum: "http://localhost:8545",
    planq: "http://localhost:8545",
    bsc: "http://localhost:8546",
  },
};

export const PHYLAX_RPC_HOSTS = [
  "https://mainnet-api.deltaswap.io",
];
