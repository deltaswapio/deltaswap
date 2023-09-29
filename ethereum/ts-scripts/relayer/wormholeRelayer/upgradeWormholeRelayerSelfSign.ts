import { deployDeltaswapRelayerImplementation } from "../helpers/deployments";
import {
  init,
  ChainInfo,
  getDeltaswapRelayer,
  writeOutputFiles,
  getOperatingChains,
} from "../helpers/env";
import { createDeltaswapRelayerUpgradeVAA } from "../helpers/vaa";

const processName = "upgradeDeltaswapRelayerSelfSign";
init();
const chains = getOperatingChains();

async function run() {
  console.log("Start!");
  const output: any = {
    wormholeRelayerImplementations: []
  };

  for (const chain of chains) {
    const coreRelayerImplementation = await deployDeltaswapRelayerImplementation(
      chain
    );
    await upgradeDeltaswapRelayer(chain, coreRelayerImplementation.address);

    output.wormholeRelayerImplementations.push(coreRelayerImplementation);
  }

  writeOutputFiles(output, processName);
}

async function upgradeDeltaswapRelayer(
  chain: ChainInfo,
  newImplementationAddress: string
) {
  console.log("upgradeDeltaswapRelayer " + chain.chainId);

  const coreRelayer = await getDeltaswapRelayer(chain);

  const tx = await coreRelayer.submitContractUpgrade(
    createDeltaswapRelayerUpgradeVAA(chain, newImplementationAddress)
  );

  await tx.wait();

  console.log(
    "Successfully upgraded the core relayer contract on " + chain.chainId
  );
}

run().then(() => console.log("Done! " + processName));
