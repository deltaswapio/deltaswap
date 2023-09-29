import {
  init,
  loadChains,
  ChainInfo,
  getDeltaswapRelayer,
  getOperatingChains,
} from "../helpers/env";
import { wait } from "../helpers/utils";
import { createRegisterChainVAA } from "../helpers/vaa";

const processName = "registerChainsDeltaswapRelayerSelfSign";
init();
const operatingChains = getOperatingChains();
const chains = loadChains();

async function run() {
  console.log("Start! " + processName);

  for (const operatingChain of operatingChains) {
    await registerChainsDeltaswapRelayer(operatingChain);
    await registerOnExistingChainsDeltaswapRelayer(operatingChain);
  }
}

async function registerChainsDeltaswapRelayer(chain: ChainInfo) {
  console.log(
    "Registering all the wormhole relayers onto Wormhole Relayer " +
      chain.chainId
  );

  const coreRelayer = await getDeltaswapRelayer(chain);
  for (const targetChain of chains) {
    try {
      await coreRelayer
        .registerDeltaswapRelayerContract(createRegisterChainVAA(targetChain))
        .then(wait);
    } catch (e) {
      console.log(
        `Error in registering chain ${targetChain.chainId} onto ${chain.chainId}`
      );
    }
  }

  console.log(
    "Did all contract registrations for the core relayer on " + chain.chainId
  );
}

async function registerOnExistingChainsDeltaswapRelayer(chain: ChainInfo) {
  console.log(
    "Registering Wormhole Relayer " +
      chain.chainId +
      " onto all the wormhole relayers"
  );
  const operatingChainIds = operatingChains.map((c) => c.chainId);
  for (const targetChain of chains) {
    if (operatingChainIds.find((x) => x === targetChain.chainId)) {
      continue;
    }
    const coreRelayer = await getDeltaswapRelayer(targetChain);
    try {
      await coreRelayer
        .registerDeltaswapRelayerContract(createRegisterChainVAA(chain))
        .then(wait);
    } catch (e) {
      console.log(
        `Error in registering chain ${chain.chainId} onto ${targetChain.chainId}`
      );
      if (targetChain.chainId === 5) {
        console.log(e);
      }
    }
  }

  console.log(
    "Did all contract registrations of the core relayer on " +
      chain.chainId +
      " onto the existing (non operating) chains"
  );
}

run().then(() => console.log("Done! " + processName));
