import {
  buildOverrides,
  deployDeltaswapRelayerImplementation,
} from "../helpers/deployments";
import {
  init,
  ChainInfo,
  getDeltaswapRelayer,
  writeOutputFiles,
  getOperatingChains,
  Deployment,
} from "../helpers/env";
import { createDeltaswapRelayerUpgradeVAA } from "../helpers/vaa";

const processName = "upgradeDeltaswapRelayerSelfSign";
init();
const operatingChains = getOperatingChains();

interface DeltaswapRelayerUpgrade {
  deltaswapRelayerImplementations: Deployment[];
}

async function run() {
  console.log("Start!");
  const output: DeltaswapRelayerUpgrade = {
    deltaswapRelayerImplementations: [],
  };

  const tasks = await Promise.allSettled(
    operatingChains.map(async (chain) => {
      const implementation = await deployDeltaswapRelayerImplementation(chain);
      await upgradeDeltaswapRelayer(chain, implementation.address);

      return implementation;
    }),
  );

  for (const task of tasks) {
    if (task.status === "rejected") {
      console.log(`DeltaswapRelayer upgrade failed. ${task.reason?.stack || task.reason}`);
    } else {
      output.deltaswapRelayerImplementations.push(task.value);
    }
  }

  writeOutputFiles(output, processName);
}

async function upgradeDeltaswapRelayer(
  chain: ChainInfo,
  newImplementationAddress: string,
) {
  console.log("upgradeDeltaswapRelayer " + chain.chainId);

  const deltaswapRelayer = await getDeltaswapRelayer(chain);

  const vaa = createDeltaswapRelayerUpgradeVAA(chain, newImplementationAddress);

  const overrides = await buildOverrides(
    () => deltaswapRelayer.estimateGas.submitContractUpgrade(vaa),
    chain,
  );
  const tx = await deltaswapRelayer.submitContractUpgrade(vaa, overrides);

  const receipt = await tx.wait();

  if (receipt.status !== 1) {
    throw new Error(
      `Failed to upgrade on chain ${chain.chainId}, tx id: ${tx.hash}`,
    );
  }
  console.log(
    "Successfully upgraded the core relayer contract on " + chain.chainId,
  );
}

run().then(() => console.log("Done! " + processName));
