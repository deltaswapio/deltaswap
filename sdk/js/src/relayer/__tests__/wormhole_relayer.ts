import { describe, expect, test } from "@jest/globals";
import { ContractReceipt, ethers } from "ethers";
import {
  getNetwork,
  isCI,
  waitForRelay,
  PRIVATE_KEY,
  getPhylaxRPC,
  PHYLAX_KEYS,
  PHYLAX_SET_INDEX,
  GOVERNANCE_EMITTER_ADDRESS,
  getArbitraryBytes32,
} from "./utils/utils";
import { getAddressInfo } from "../consts";
import { getDefaultProvider } from "../relayer/helpers";
import {
  relayer,
  ethers_contracts,
  tryNativeToUint8Array,
  ChainId,
  CHAINS,
  CONTRACTS,
  ChainName,
  Network,
} from "../../../";
import { GovernanceEmitter, MockPhylaxs } from "../../../src/mock";
import { Implementation__factory } from "../../ethers-contracts";
import {manualDelivery} from "../relayer";
import { NodeHttpTransport } from "@improbable-eng/grpc-web-node-http-transport";
import { packEVMExecutionInfoV1 } from "../structs";

const network: Network = getNetwork();
const ci: boolean = isCI();

const sourceChain = network == "DEVNET" ? "ethereum" : "celo";
const targetChain = network == "DEVNET" ? "bsc" : "avalanche";

const testIfDevnet = () => (network == "DEVNET" ? test : test.skip);
const testIfNotDevnet = () => (network != "DEVNET" ? test : test.skip);

type TestChain = {
  chainId: ChainId;
  name: ChainName;
  provider: ethers.providers.Provider;
  wallet: ethers.Wallet;
  deltaswapRelayerAddress: string;
  mockIntegrationAddress: string;
  deltaswapRelayer: ethers_contracts.DeltaswapRelayer;
  mockIntegration: ethers_contracts.MockRelayerIntegration;
};

const createTestChain = (name: ChainName) => {
  const provider = getDefaultProvider(network, name, ci);
  const addressInfo = getAddressInfo(name, network);
  if (process.env.DEV) {
    // Via ir is off -> different deltaswap relayer address
    addressInfo.deltaswapRelayerAddress =
      "0x53855d4b64E9A3CF59A84bc768adA716B5536BC5";
  }
  if (network == "MAINNET")
    addressInfo.mockIntegrationAddress =
      "0xa507Ff8D183D2BEcc9Ff9F82DFeF4b074e1d0E05";
  if (network == "MAINNET")
    addressInfo.mockDeliveryProviderAddress =
      "0x7A0a53847776f7e94Cc35742971aCb2217b0Db81";

  if (!addressInfo.deltaswapRelayerAddress)
    throw Error(`No core relayer address for ${name}`);
  if (!addressInfo.mockIntegrationAddress)
    throw Error(`No mock relayer integration address for ${name}`);
  const wallet = new ethers.Wallet(PRIVATE_KEY, provider);
  const deltaswapRelayer = ethers_contracts.DeltaswapRelayer__factory.connect(
    addressInfo.deltaswapRelayerAddress,
    wallet
  );
  const mockIntegration =
    ethers_contracts.MockRelayerIntegration__factory.connect(
      addressInfo.mockIntegrationAddress,
      wallet
    );
  const result: TestChain = {
    chainId: CHAINS[name],
    name,
    provider,
    wallet,
    deltaswapRelayerAddress: addressInfo.deltaswapRelayerAddress,
    mockIntegrationAddress: addressInfo.mockIntegrationAddress,
    deltaswapRelayer,
    mockIntegration,
  };
  return result;
};

const source = createTestChain(sourceChain);
const target = createTestChain(targetChain);

const myMap = new Map<ChainName, ethers.providers.Provider>();
myMap.set(sourceChain, source.provider);
myMap.set(targetChain, target.provider);
const optionalParams = {
  environment: network,
  sourceChainProvider: source.provider,
  targetChainProviders: myMap,
  deltaswapRelayerAddress: source.deltaswapRelayerAddress,
};
const optionalParamsTarget = {
  environment: network,
  sourceChainProvider: target.provider,
  targetChainProviders: myMap,
  deltaswapRelayerAddress: target.deltaswapRelayerAddress,
};

// for signing deltaswap messages
const phylaxs = new MockPhylaxs(PHYLAX_SET_INDEX, PHYLAX_KEYS);

// for generating governance deltaswap messages
const governance = new GovernanceEmitter(GOVERNANCE_EMITTER_ADDRESS);

const phylaxIndices = process.env.NUM_PHYLAXS
    ? [...Array(parseInt(process.env.NUM_PHYLAXS)).keys()]
    : ci
        ? [0, 1]
        : [0];

const REASONABLE_GAS_LIMIT = 500000;
const TOO_LOW_GAS_LIMIT = 10000;

const deltaswapRelayerAddresses = new Map<ChainName, string>();
deltaswapRelayerAddresses.set(sourceChain, source.deltaswapRelayerAddress);
deltaswapRelayerAddresses.set(targetChain, target.deltaswapRelayerAddress);

const getStatus = async (
  txHash: string,
  _sourceChain?: ChainName,
  index?: number
): Promise<string> => {
  const info = (await relayer.getDeltaswapRelayerInfo(
    _sourceChain || sourceChain,
    txHash,
    {
      environment: network,
      targetChainProviders: myMap,
      sourceChainProvider: myMap.get(_sourceChain || sourceChain),
      deltaswapRelayerAddresses,
    }
  )) as relayer.DeliveryInfo;
  return info.targetChainStatus.events[index ? index : 0].status;
};

const testSend = async (
  payload: string,
  sendToSourceChain?: boolean,
  notEnoughValue?: boolean
): Promise<ContractReceipt> => {
  const value = await relayer.getPrice(
    sourceChain,
    sendToSourceChain ? sourceChain : targetChain,
    notEnoughValue ? TOO_LOW_GAS_LIMIT : REASONABLE_GAS_LIMIT,
    optionalParams
  );
  console.log(`Quoted gas delivery fee: ${value}`);
  const tx = await source.mockIntegration.sendMessage(
    payload,
    sendToSourceChain ? source.chainId : target.chainId,
    notEnoughValue ? TOO_LOW_GAS_LIMIT : REASONABLE_GAS_LIMIT,
    0,
    { value, gasLimit: REASONABLE_GAS_LIMIT }
  );
  console.log(`Sent delivery request! Transaction hash ${tx.hash}`);
  await tx.wait();
  console.log("Message confirmed!");

  return tx.wait();
};

describe("Deltaswap Relayer Tests", () => {
  test("Executes a Delivery Success", async () => {
    const arbitraryPayload = getArbitraryBytes32();
    console.log(`Sent message: ${arbitraryPayload}`);

    const rx = await testSend(arbitraryPayload);

    await waitForRelay();

    console.log("Checking if message was relayed");
    const message = await target.mockIntegration.getMessage();
    expect(message).toBe(arbitraryPayload);
  });

  test("Executes a Delivery Success With Additional VAAs", async () => {
    const arbitraryPayload = getArbitraryBytes32();
    console.log(`Sent message: ${arbitraryPayload}`);

    const deltaswap = Implementation__factory.connect(
      CONTRACTS[network][sourceChain].core || "",
      source.wallet
    );
    const deliverySeq = await deltaswap.nextSequence(source.wallet.address);
    const msgTx = await deltaswap.publishMessage(0, arbitraryPayload, 200);
    await msgTx.wait();

    const value = await relayer.getPrice(
      sourceChain,
      targetChain,
      REASONABLE_GAS_LIMIT * 2,
      optionalParams
    );
    console.log(`Quoted gas delivery fee: ${value}`);

    const tx = await source.mockIntegration.sendMessageWithAdditionalVaas(
      [],
      target.chainId,
      REASONABLE_GAS_LIMIT * 2,
      0,
      [
        relayer.createVaaKey(
          source.chainId,
          Buffer.from(tryNativeToUint8Array(source.wallet.address, "ethereum")),
          deliverySeq
        ),
      ],
      { value }
    );

    console.log(`Sent tx hash: ${tx.hash}`);

    const rx = await tx.wait();

    await waitForRelay();

    console.log("Checking if message was relayed");
    const message = (await target.mockIntegration.getDeliveryData())
      .additionalVaas[0];
    const parsedMessage = await deltaswap.parseVM(message);
    expect(parsedMessage.payload).toBe(arbitraryPayload);
  });

    testIfNotDevnet()(
        "Executes a Delivery Success with manual delivery",
        async () => {
            const arbitraryPayload = getArbitraryBytes32();
            console.log(`Sent message: ${arbitraryPayload}`);

            const deliverySeq = await Implementation__factory.connect(
                CONTRACTS[network][sourceChain].core || "",
                source.provider
            ).nextSequence(source.deltaswapRelayerAddress);

            const rx = await testSend(arbitraryPayload, false, true);

            await waitForRelay();

            // confirm that the message was not relayed successfully
            {
                const message = await target.mockIntegration.getMessage();
                expect(message).not.toBe(arbitraryPayload);
            }
            const [value, refundPerGasUnused] = await relayer.getPriceAndRefundInfo(
                sourceChain,
                targetChain,
                REASONABLE_GAS_LIMIT,
                optionalParams
            );

            const priceInfo = await manualDelivery(
                sourceChain,
                rx.transactionHash,
                {deltaswapRelayerAddresses, ...optionalParams},
                true,
                {
                    newExecutionInfo: Buffer.from(
                        packEVMExecutionInfoV1({
                            gasLimit: ethers.BigNumber.from(REASONABLE_GAS_LIMIT),
                            targetChainRefundPerGasUnused:
                                ethers.BigNumber.from(refundPerGasUnused),
                        }).substring(2),
                        "hex"
                    ),
                    newReceiverValue: ethers.BigNumber.from(0),
                    redeliveryHash: Buffer.from(
                        ethers.utils.keccak256("0x1234").substring(2),
                        "hex"
                    ), // fake a redelivery
                }
            );

            console.log(`Price: ${priceInfo.quote} of ${priceInfo.targetChain} wei`);

            const deliveryRx = await manualDelivery(
                sourceChain,
                rx.transactionHash,
                {deltaswapRelayerAddresses, ...optionalParams},
                false,
                {
                    newExecutionInfo: Buffer.from(
                        packEVMExecutionInfoV1({
                            gasLimit: ethers.BigNumber.from(REASONABLE_GAS_LIMIT),
                            targetChainRefundPerGasUnused:
                                ethers.BigNumber.from(refundPerGasUnused),
                        }).substring(2),
                        "hex"
                    ),
                    newReceiverValue: ethers.BigNumber.from(0),
                    redeliveryHash: Buffer.from(
                        ethers.utils.keccak256("0x1234").substring(2),
                        "hex"
                    ), // fake a redelivery
                },
                target.wallet
            );
            console.log("Manual delivery tx hash", deliveryRx.txHash);

            console.log("Checking if message was relayed");
            const message = await target.mockIntegration.getMessage();
            expect(message).toBe(arbitraryPayload);
        }
    );

  testIfDevnet()("Test getPrice in Typescript SDK", async () => {
    const price = await relayer.getPrice(
      sourceChain,
      targetChain,
      200000,
      optionalParams
    );
    expect(price.toString()).toBe("165000000000000000");
  });

  test("Executes a delivery with a Cross Chain Refund", async () => {
    const arbitraryPayload = getArbitraryBytes32();
    console.log(`Sent message: ${arbitraryPayload}`);
    const value = await relayer.getPrice(
      sourceChain,
      targetChain,
      REASONABLE_GAS_LIMIT,
      optionalParams
    );
    console.log(`Quoted gas delivery fee: ${value}`);
    const startingBalance = await source.wallet.getBalance();

    const tx = await relayer.sendToEvm(
      source.wallet,
      sourceChain,
      targetChain,
      target.deltaswapRelayerAddress, // This is an address that exists but doesn't implement the IDeltaswap interface, so should result in Receiver Failure
      Buffer.from("hi!"),
      REASONABLE_GAS_LIMIT,
      { value, gasLimit: REASONABLE_GAS_LIMIT },
      optionalParams
    );
    console.log("Sent delivery request!");
    await tx.wait();
    console.log("Message confirmed!");
    const endingBalance = await source.wallet.getBalance();

    await waitForRelay();

    const info = (await relayer.getDeltaswapRelayerInfo(sourceChain, tx.hash, {
      deltaswapRelayerAddresses,
      ...optionalParams,
    })) as relayer.DeliveryInfo;

    await waitForRelay();

    const newEndingBalance = await source.wallet.getBalance();

    console.log(`Quoted gas delivery fee: ${value}`);
    console.log(
      `Cost (including gas) ${startingBalance.sub(endingBalance).toString()}`
    );
    const refund = newEndingBalance.sub(endingBalance);
    console.log(`Refund: ${refund.toString()}`);
    console.log(
      `As a percentage of original value: ${newEndingBalance
        .sub(endingBalance)
        .mul(100)
        .div(value)
        .toString()}%`
    );
    console.log("Confirming refund is nonzero");
    expect(refund.gt(0)).toBe(true);
  });

  test("Executes a Receiver Failure", async () => {
    const arbitraryPayload = getArbitraryBytes32();
    console.log(`Sent message: ${arbitraryPayload}`);

    const rx = await testSend(arbitraryPayload, false, true);

    await waitForRelay();

    const message = await target.mockIntegration.getMessage();
    expect(message).not.toBe(arbitraryPayload);
  });

  test("Executes a receiver failure and then redelivery through SDK", async () => {
    const arbitraryPayload = getArbitraryBytes32();
    console.log(`Sent message: ${arbitraryPayload}`);

    const rx = await testSend(arbitraryPayload, false, true);

    await waitForRelay();

    const message = await target.mockIntegration.getMessage();
    expect(message).not.toBe(arbitraryPayload);

    const value = await relayer.getPrice(
      sourceChain,
      targetChain,
      REASONABLE_GAS_LIMIT,
      optionalParams
    );

    const info = (await relayer.getDeltaswapRelayerInfo(
      sourceChain,
      rx.transactionHash,
      { deltaswapRelayerAddresses, ...optionalParams }
    )) as relayer.DeliveryInfo;

    console.log("Redelivering message");
    const redeliveryReceipt = await relayer.resend(
      source.wallet,
      sourceChain,
      targetChain,
      network,
      relayer.createVaaKey(
        source.chainId,
        Buffer.from(
          tryNativeToUint8Array(source.deltaswapRelayerAddress, "ethereum")
        ),
        info.sourceDeliverySequenceNumber
      ),
      REASONABLE_GAS_LIMIT,
      0,
      await source.deltaswapRelayer.getDefaultDeliveryProvider(),
      [getPhylaxRPC(network, ci)],
      {
        value: value,
        gasLimit: REASONABLE_GAS_LIMIT,
      },
      { transport: NodeHttpTransport() },
      { deltaswapRelayerAddress: source.deltaswapRelayerAddress }
    );

    console.log("redelivery tx:", redeliveryReceipt.hash);

    await redeliveryReceipt.wait();

    await waitForRelay();

    console.log("Checking if message was relayed after redelivery");
    const message2 = await target.mockIntegration.getMessage();
    expect(message2).toBe(arbitraryPayload);

    //Can extend this to look for redelivery event
  });

  // GOVERNANCE TESTS

  testIfDevnet()("Governance: Test Registering Chain", async () => {
    const chain = 24;

    const currentAddress =
      await source.deltaswapRelayer.getRegisteredDeltaswapRelayerContract(chain);
    console.log(
      `For Chain ${source.chainId}, registered chain ${chain} address: ${currentAddress}`
    );

    const expectedNewRegisteredAddress =
      "0x0000000000000000000000001234567890123456789012345678901234567892";

    const timestamp = (await source.wallet.provider.getBlock("latest"))
      .timestamp;

    const firstMessage = governance.publishDeltaswapRelayerRegisterChain(
      timestamp,
      chain,
      expectedNewRegisteredAddress
    );
    const firstSignedVaa = phylaxs.addSignatures(
      firstMessage,
      phylaxIndices
    );

    let tx = await source.deltaswapRelayer.registerDeltaswapRelayerContract(
      firstSignedVaa,
      { gasLimit: REASONABLE_GAS_LIMIT }
    );
    await tx.wait();

    const newRegisteredAddress =
      await source.deltaswapRelayer.getRegisteredDeltaswapRelayerContract(chain);

    expect(newRegisteredAddress).toBe(expectedNewRegisteredAddress);
  });

  testIfDevnet()(
    "Governance: Test Setting Default Relay Provider",
    async () => {
      const currentAddress =
        await source.deltaswapRelayer.getDefaultDeliveryProvider();
      console.log(
        `For Chain ${source.chainId}, default relay provider: ${currentAddress}`
      );

      const expectedNewDefaultDeliveryProvider =
        "0x1234567890123456789012345678901234567892";

      const timestamp = (await source.wallet.provider.getBlock("latest"))
        .timestamp;
      const chain = source.chainId;
      const firstMessage =
        governance.publishDeltaswapRelayerSetDefaultDeliveryProvider(
          timestamp,
          chain,
          expectedNewDefaultDeliveryProvider
        );
      const firstSignedVaa = phylaxs.addSignatures(
        firstMessage,
        phylaxIndices
      );

      let tx = await source.deltaswapRelayer.setDefaultDeliveryProvider(
        firstSignedVaa
      );
      await tx.wait();

      const newDefaultDeliveryProvider =
        await source.deltaswapRelayer.getDefaultDeliveryProvider();

      expect(newDefaultDeliveryProvider).toBe(
        expectedNewDefaultDeliveryProvider
      );

      const inverseFirstMessage =
        governance.publishDeltaswapRelayerSetDefaultDeliveryProvider(
          timestamp,
          chain,
          currentAddress
        );
      const inverseFirstSignedVaa = phylaxs.addSignatures(
        inverseFirstMessage,
        phylaxIndices
      );

      tx = await source.deltaswapRelayer.setDefaultDeliveryProvider(
        inverseFirstSignedVaa
      );
      await tx.wait();

      const originalDefaultDeliveryProvider =
        await source.deltaswapRelayer.getDefaultDeliveryProvider();

      expect(originalDefaultDeliveryProvider).toBe(currentAddress);
    }
  );

  testIfDevnet()("Governance: Test Upgrading Contract", async () => {
    const IMPLEMENTATION_STORAGE_SLOT =
      "0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc";

    const getImplementationAddress = () =>
      source.provider.getStorageAt(
        source.deltaswapRelayer.address,
        IMPLEMENTATION_STORAGE_SLOT
      );

    console.log(
      `Current Implementation address: ${await getImplementationAddress()}`
    );

    const deltaswapAddress = CONTRACTS[network][sourceChain].core || "";

    const newDeltaswapRelayerImplementationAddress = (
      await new ethers_contracts.DeltaswapRelayer__factory(source.wallet)
        .deploy(deltaswapAddress)
        .then((x) => x.deployed())
    ).address;

    console.log(`Deployed!`);
    console.log(
      `New core relayer implementation: ${newDeltaswapRelayerImplementationAddress}`
    );

    const timestamp = (await source.wallet.provider.getBlock("latest"))
      .timestamp;
    const chain = source.chainId;
    const firstMessage = governance.publishDeltaswapRelayerUpgradeContract(
      timestamp,
      chain,
      newDeltaswapRelayerImplementationAddress
    );
    const firstSignedVaa = phylaxs.addSignatures(
      firstMessage,
      phylaxIndices
    );

    let tx = await source.deltaswapRelayer.submitContractUpgrade(firstSignedVaa);

    expect(
      ethers.utils.getAddress((await getImplementationAddress()).substring(26))
    ).toBe(ethers.utils.getAddress(newDeltaswapRelayerImplementationAddress));
  });

    testIfNotDevnet()("Checks the status of a message", async () => {
        const txHash =
            "0xa75e4100240e9b498a48fa29de32c9e62ec241bf4071a3c93fde0df5de53c507";
        const mySourceChain: ChainName = "celo";
        const environment: Network = "TESTNET";

        const info = await relayer.getDeltaswapRelayerInfo(mySourceChain, txHash, {
            environment,
        });
        console.log(info.stringified);
    });

    testIfNotDevnet()("Tests custom manual delivery", async () => {
        const txHash =
            "0xc57d12cc789e4e9fa50d496cea62c2a0f11a7557c8adf42b3420e0585ba1f911";
        const mySourceChain: ChainName = "arbitrum";
        const targetProvider = undefined;
        const environment: Network = "TESTNET";

        const info = await relayer.getDeltaswapRelayerInfo(mySourceChain, txHash, {
            environment,
        });
        console.log(info.stringified);

        const priceInfo = await manualDelivery(
            mySourceChain,
            txHash,
            {environment},
            true
        );
        console.log(`Price info: ${JSON.stringify(priceInfo)}`);

        const signer = new ethers.Wallet(
            PRIVATE_KEY,
            targetProvider
                ? new ethers.providers.JsonRpcProvider(targetProvider)
                : getDefaultProvider(environment, priceInfo.targetChain)
        );

        console.log(
            `Price: ${ethers.utils.formatEther(priceInfo.quote)} of ${
                priceInfo.targetChain
            } currency`
        );
        const balance = await signer.getBalance();
        console.log(
            `My balance: ${ethers.utils.formatEther(balance)} of ${
                priceInfo.targetChain
            } currency`
        );

        const deliveryRx = await manualDelivery(
            mySourceChain,
            txHash,
            {environment},
            false,
            undefined,
            signer
        );
        console.log("Manual delivery tx hash", deliveryRx.txHash);
    });
});

function sleep(ms: number): Promise<void> {
  return new Promise((r) => setTimeout(() => r(), ms));
}
