import {
  deployDeltaswapRelayerImplementation,
  deployDeltaswapRelayerProxy,
} from "../helpers/deployments";
import {
  init,
  writeOutputFiles,
  getDeliveryProviderAddress,
  Deployment,
  getOperationDescriptor,
  loadLastRun,
} from "../helpers/env";

const processName = "deployDeltaswapRelayer";
init();
const operation = getOperationDescriptor();

interface DeltaswapRelayerDeployment {
  deltaswapRelayerImplementations: Deployment[];
  deltaswapRelayerProxies: Deployment[];
}

async function run() {
  console.log("Start! " + processName);

  const lastRun: DeltaswapRelayerDeployment | undefined =
    loadLastRun(processName);
  const deployments: DeltaswapRelayerDeployment = {
      deltaswapRelayerImplementations: lastRun?.deltaswapRelayerImplementations?.filter(isSupportedChain) || [],
      deltaswapRelayerProxies: lastRun?.deltaswapRelayerProxies?.filter(isSupportedChain) || [],
  };

  for (const chain of operation.operatingChains) {
    console.log(`Deploying for chain ${chain.chainId}...`);
    const coreRelayerImplementation = await deployDeltaswapRelayerImplementation(
      chain,
    );
    const coreRelayerProxy = await deployDeltaswapRelayerProxy(
      chain,
      coreRelayerImplementation.address,
      getDeliveryProviderAddress(chain),
    );

    deployments.deltaswapRelayerImplementations.push(
      coreRelayerImplementation,
    );
    deployments.deltaswapRelayerProxies.push(coreRelayerProxy);
    console.log("");
  }

  writeOutputFiles(deployments, processName);
}

function isSupportedChain(deploy: Deployment): boolean {
  const item = operation.supportedChains.find((chain) => {
    return deploy.chainId === chain.chainId;
  });
  return item !== undefined;
}

run().then(() => console.log("Done! " + processName));
