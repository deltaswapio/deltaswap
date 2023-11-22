import {
  deployDeltaswapRelayerImplementation,
} from "../helpers/deployments";
import {
  init,
  writeOutputFiles,
  getOperatingChains,
  Deployment,
} from "../helpers/env";

const processName = "deployDeltaswapRelayerImplementation";
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

run().then(() => console.log("Done! " + processName));
