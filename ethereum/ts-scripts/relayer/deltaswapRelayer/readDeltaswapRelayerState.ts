import {
  init,
  ChainInfo,
  getDeltaswapRelayerAddress,
  getProvider,
  writeOutputFiles,
  getDeltaswapRelayer,
  getOperatingChains,
  loadChains,
} from "../helpers/env";

const processName = "readDeltaswapRelayerContractState";
init();
const allChains = loadChains();
const operatingChains = getOperatingChains();

async function run() {
  console.log("Start! " + processName);

  const states: any = [];

  for (const chain of operatingChains) {
    const state = await readState(chain);
    if (state) {
      printState(state);
      states.push(state);
    }
  }

  writeOutputFiles(states, processName);
}

type DeltaswapRelayerContractState = {
  chainId: number;
  contractAddress: string;
  defaultProvider: string;
  registeredContracts: { chainId: number; contract: string }[];
};

async function readState(
  chain: ChainInfo
): Promise<DeltaswapRelayerContractState | null> {
  console.log(
    "Gathering core relayer contract status for chain " + chain.chainId,
  );

  try {
    const contractAddress = await getDeltaswapRelayerAddress(chain);
    console.log("Querying " + contractAddress);

    const coreRelayer = await getDeltaswapRelayer(chain, getProvider(chain));

    console.log("Querying default provider for code");
    const provider = getProvider(chain);
    const codeReceipt = await provider.getCode(contractAddress);
    console.log("Code: " + codeReceipt);

    const registeredContracts: { chainId: number; contract: string }[] = [];

    for (const chainInfo of allChains) {
      registeredContracts.push({
        chainId: chainInfo.chainId,
        contract: (
          await coreRelayer.getRegisteredDeltaswapRelayerContract(
            chainInfo.chainId
          )
        ).toString(),
      });
    }

    const defaultProvider = await coreRelayer.getDefaultDeliveryProvider();
    return {
      chainId: chain.chainId,
      contractAddress,
      defaultProvider,
      registeredContracts,
    };
  } catch (e) {
    console.error(e);
    console.log("Failed to gather status for chain " + chain.chainId);
  }

  return null;
}

function printState(state: DeltaswapRelayerContractState) {
  console.log("");
  console.log("DeltaswapRelayer: ");
  printFixed("Chain ID: ", state.chainId.toString());
  printFixed("Contract Address:", state.contractAddress);
  printFixed("Default Provider:", state.defaultProvider);

  console.log("");

  printFixed("Registered DeltaswapRelayers", "");
  state.registeredContracts.forEach((x) => {
    printFixed("  Chain: " + x.chainId, x.contract);
  });
  console.log("");
}

function printFixed(title: string, content: string) {
  const length = 80;
  const spaces = length - title.length - content.length;
  let str = "";
  if (spaces > 0) {
    for (let i = 0; i < spaces; i++) {
      str = str + " ";
    }
  }
  console.log(title + str + content);
}

run().then(() => console.log("Done! " + processName));
