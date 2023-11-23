import {
  init,
  getOperatingChains,
  getDeltaswapRelayer,
  ChainInfo,
} from "../helpers/env";
import { buildOverrides } from "../helpers/deployments";

import { extractChainToBeRegisteredFromRegisterChainVaa } from "../helpers/vaa";

import { inspect } from "util";

interface RegisterChainDescriptor {
  vaa: Buffer;
  chainToRegister: number;
}

const processName = "registerDeltaswapRelayer";
init();
const chains = getOperatingChains();

const zeroBytes32 =
  "0x0000000000000000000000000000000000000000000000000000000000000000";

/**
 * These are the registration VAAs for mainnet.
 */
const base64RegistrationVaas = [
  "AQAAAAEHADqL9QobrXc1fnh8LP+A0EbxNHkj+cOfsmhp1dpFIq/VNoGj0hIPa2Y1wxATI0ghtQG3zJc+a11+7k9SRpUZwXoAAhHcBLIuQIyp0SkvFQuinab+H5sH5UaEAtQzDD6CJAslKdy3wGjzPMiJEO0iH7IXSfMop+eq7lZNUNC1pbzTBQABBJHQ3BjzLq+38tPPXapWhnTRLre/z6uG0bGbqf7XWwv2aa33Ptt4fXtJxvrGLkucehLR7k1kneOvaDsrlzjUi2EBBUA5DIlqSE0voke4NC+94t0S/VWiQxgKwiFLutNrGuCQXiAWTXixKrd/vESX8ygdiTds6aRKjOxvs5LM7YfTonYABov/gpRsje8vs+o2GXLUNLC0/r9bsqwmzFd6rvAGi63Xb5QPCFGC1DWkUhZegSL7tYTBhXoFGf9VG2kWZtLXXXsBB37rpvgpWpOKbDvftp86csp7YT9zZA+RJUXA8bhwqVFgReuy/6QGYJtR50UPQ6QfMjZeLBbvarwZLA9xjDvdFiYBCIEiEVjVsveNyhEzIONJ/9jV5mdPJn+nPBSDO2Lbd41vG/Az2EE1FgOosQw19mu6ZWbDSQML3NCyn0W2gncsd+wBAAAAAJ0JwxkAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE72SCFLmd03EgAAAAAAAAAAAAAAAAAAAAAABXb3JtaG9sZVJlbGF5ZXIBAAAbngAAAAAAAAAAAAAAAOOLvm7/VMYPD/OtMPXEKfYzsRfG",
  "AQAAAAEHANQ+54wplFAes9Bs4muWnnq4EU5MFvQzJFuUM+t+EYL2Xsd5VlRNI1ccr+1A/pAprKEKIv5einMQAUyOgkqSRKgAAkqGfpykzL+p/c9Zv4J0jypVmQIFvMc1nsDjH6d5K/sSQb3q0MVbIJxuSh/Hp29/ZErKIRl/K0KdObtwKScRaIUBBNfJ0ptV0qS3Y4z3oB2ivNPMkA/ZzC5ljR6DyRwc1WTiWj29SFXE45KQNsRlCeqhrwyu0mkVrpJ3140yp4G26FgABRTpr1sl9QmMoo8j3avxLJz08P/bQbKrT7TkGeFk2s3DaujhpuPP3Jy9w9Sli6AGhhdP5EAIKk6IrZQmguKKBKgBBhwvpbiERM+bzYcN1D8mhxbf7P/KVkjh3b7TyH6Sj2ibJAVbQtuMjaSQs8rY7l496xL0apMgI6YNKJEqe2oazcoBB8zbHxQezbZh9wHNYtwqkpMs4+cKwYA0H6vgJx8djC8rVFwX5N6La6xPhaaBJ2Z5s+h1lRftl3jiomvn/O8nqMYACKts63v6+wjI46lOW7iVwh7Cdvf4pGydwB4BfYbFg20uEkDbgL+1CUDpW2w7EFT0TuDB+m0fLSyLYQnXdqaqFq8AAAAAACp/748AAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEKlGTaFwHFgMgAAAAAAAAAAAAAAAAAAAAAABXb3JtaG9sZVJlbGF5ZXIBAAAABAAAAAAAAAAAAAAAAGXHGSswF7xPHjCmyPbYijIcMTgU",
];

async function run() {
  console.log(`Start! ${processName}`);

  const registrationVaas = base64RegistrationVaas.map((vaa) => {
    const vaaBuf = Buffer.from(vaa, "base64");
    return {
      vaa: vaaBuf,
      chainToRegister: extractChainToBeRegisteredFromRegisterChainVaa(vaaBuf),
    };
  });
  const tasks = await Promise.allSettled(
    chains.map((chain) => register(chain, registrationVaas)),
  );

  for (const task of tasks) {
    if (task.status === "rejected") {
      console.error(`Register failed: ${inspect(task.reason)}`);
    } else {
      console.log(`Register succeeded for chain ${task.value}`);
    }
  }
}

async function register(
  chain: ChainInfo,
  registrationVaas: RegisterChainDescriptor[],
) {
  console.log(`Registering DeltaswapRelayers in chain ${chain.chainId}...`);
  const deltaswapRelayer = await getDeltaswapRelayer(chain);

  // TODO: check for already registered VAAs
  for (const { vaa, chainToRegister } of registrationVaas) {
    const registrationAddress =
      await deltaswapRelayer.getRegisteredDeltaswapRelayerContract(
        chainToRegister,
      );
    if (registrationAddress !== zeroBytes32) {
      // We skip chains that are already registered.
      // Note that reregistrations aren't allowed.
      continue;
    }

    const overrides = await buildOverrides(
      () => deltaswapRelayer.estimateGas.registerDeltaswapRelayerContract(vaa),
      chain,
    );
    const tx = await deltaswapRelayer.registerDeltaswapRelayerContract(
      vaa,
      overrides,
    );
    const receipt = await tx.wait();

    if (receipt.status !== 1) {
      throw new Error(
        `Failed registration for chain ${chain.chainId}, tx ${tx.hash}`,
      );
    }
  }

  return chain.chainId;
}

run().then(() => console.log("Done! " + processName));
