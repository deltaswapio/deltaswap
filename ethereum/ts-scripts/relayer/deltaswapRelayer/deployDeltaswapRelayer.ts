import {
  deployDeltaswapRelayerImplementation,
  deployDeltaswapRelayerProxy,
} from "../helpers/deployments";
import {
  init,
  writeOutputFiles,
  getDeliveryProviderAddress,
  getOperatingChains,
  Deployment,
} from "../helpers/env";

const processName = "deployDeltaswapRelayer";
init();
const chains = getOperatingChains();

async function run() {
  console.log("Start! " + processName);

  const output: Record<string, Deployment[]> = {
    deltaswapRelayerImplementations: [],
    deltaswapRelayerProxies: [],
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

    output.deltaswapRelayerImplementations.push(coreRelayerImplementation);
    output.deltaswapRelayerProxies.push(coreRelayerProxy);
    console.log("");
  }

  writeOutputFiles(output, processName);
}

run().then(() => console.log("Done! " + processName));
