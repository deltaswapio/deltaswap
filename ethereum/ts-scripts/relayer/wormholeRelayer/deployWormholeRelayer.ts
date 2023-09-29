import {
  deployDeltaswapRelayerImplementation,
  deployDeltaswapRelayerProxy,
} from "../helpers/deployments";
import {
  init,
  writeOutputFiles,
  getDeliveryProviderAddress,
  getOperatingChains,
} from "../helpers/env";

const processName = "deployDeltaswapRelayer";
init();
const chains = getOperatingChains();

async function run() {
  console.log("Start! " + processName);

  const output: any = {
    wormholeRelayerImplementations: [],
    wormholeRelayerProxies: [],
  };

  for (const chain of chains) {
    console.log(`Deploying for chain ${chain.chainId}...`);
    const coreRelayerImplementation = await deployDeltaswapRelayerImplementation(
      chain
    );
    const coreRelayerProxy = await deployDeltaswapRelayerProxy(
      chain,
      coreRelayerImplementation.address,
      getDeliveryProviderAddress(chain)
    );

    output.wormholeRelayerImplementations.push(coreRelayerImplementation);
    output.wormholeRelayerProxies.push(coreRelayerProxy);
    console.log("");
  }

  writeOutputFiles(output, processName);
}

run().then(() => console.log("Done! " + processName));
