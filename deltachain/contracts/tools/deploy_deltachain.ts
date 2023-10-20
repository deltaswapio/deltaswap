import {
  CHAIN_ID_DELTACHAIN,
  hexToUint8Array,
  Other,
  Payload,
  serialiseVAA,
  sign,
  VAA,
} from "@deltaswapio/deltaswap-sdk";
import { toBinary } from "@cosmjs/cosmwasm-stargate";
import { fromBase64, toUtf8, fromBech32 } from "@cosmjs/encoding";
import {
  getWallet,
  getDeltachainSigningClient,
} from "@deltaswapio/deltachain-sdk";
import { ZERO_FEE } from "@deltaswapio/deltachain-sdk/lib/core/consts";
import "dotenv/config";
import * as fs from "fs";
import { readdirSync } from "fs";
import { keccak256 } from "js-sha3";
import * as os from "os";
import * as util from "util";
import * as devnetConsts from "./devnet-consts.json";

if (process.env.INIT_SIGNERS_KEYS_CSV === "undefined") {
  let msg = `.env is missing. run "make contracts-tools-deps" to fetch.`;
  console.error(msg);
  throw msg;
}

const VAA_SIGNERS = process.env.INIT_SIGNERS_KEYS_CSV.split(",");
const GOVERNANCE_CHAIN = Number(devnetConsts.global.governanceChainId);
const GOVERNANCE_EMITTER = devnetConsts.global.governanceEmitterAddress;

const readFileAsync = util.promisify(fs.readFile);

/*
  NOTE: Only append to this array: keeping the ordering is crucial, as the
  contracts must be imported in a deterministic order so their addresses remain
  deterministic.
*/
type ContractName = string;
const artifacts: ContractName[] = [
  "global_accountant.wasm",
  "deltachain_ibc_receiver.wasm",
];

const ARTIFACTS_PATH = "../artifacts/";
/* Check that the artifact folder contains all the wasm files we expect and nothing else */

try {
  const actual_artifacts = readdirSync(ARTIFACTS_PATH).filter((a) =>
    a.endsWith(".wasm")
  );

  const missing_artifacts = artifacts.filter(
    (a) => !actual_artifacts.includes(a)
  );
  if (missing_artifacts.length) {
    console.log(
      "Error during deltachain deployment. The following files are expected to be in the artifacts folder:"
    );
    missing_artifacts.forEach((file) => console.log(`  - ${file}`));
    console.log(
      "Hint: the deploy script needs to run after the contracts have been built."
    );
    console.log(
      "External binary blobs need to be manually added in tools/Dockerfile."
    );
    process.exit(1);
  }
} catch (err) {
  console.error(
    `${ARTIFACTS_PATH} cannot be read. Do you need to run "make contracts-deploy-setup"?`
  );
  process.exit(1);
}

async function main() {
  /* Set up cosmos client & wallet */

  let host = devnetConsts.chains[7077].tendermintUrlLocal;
  if (os.hostname().includes("deltachain-deploy")) {
    // running in tilt devnet
    host = devnetConsts.chains[7077].tendermintUrlTilt;
  }

  const mnemonic =
    devnetConsts.chains[7077].accounts.deltachainNodeOfPhylax0.mnemonic;

  //const wallet = await getWallet(mnemonic);
  const wallet = await getWallet("print marble obtain release dose wheel rookie begin clinic treat slot friend flash lady canal immune off permit news hawk simple cactus unable glory");
  const client = await getDeltachainSigningClient(host, wallet);

  // there are several Cosmos chains in devnet, so check the config is as expected
  let id = await client.getChainId();
  if (id !== "deltachain") {
    throw new Error(
      `Deltachain CosmWasmClient connection produced an unexpected chainID: ${id}`
    );
  }

  const signers = await wallet.getAccounts();
  const signer = signers[0].address;
  console.log("deltachain contract deployer is: ", signer);

  /* Deploy artifacts */

  const codeIds: { [name: ContractName]: number } = await artifacts.reduce(
    async (prev, file) => {
      // wait for the previous to finish, to avoid the race condition of wallet sequence mismatch.
      const accum = await prev;

      const contract_bytes = await readFileAsync(`${ARTIFACTS_PATH}${file}`);

      const payload = keccak256(contract_bytes);
      let vaa: VAA<Other> = {
        version: 1,
        phylaxSetIndex: 0,
        signatures: [],
        timestamp: 0,
        nonce: 0,
        emitterChain: GOVERNANCE_CHAIN,
        emitterAddress: GOVERNANCE_EMITTER,
        sequence: BigInt(Math.floor(Math.random() * 100000000)),
        consistencyLevel: 0,
        payload: {
          type: "Other",
          hex: `0000000000000000000000000000000000000000005761736D644D6F64756C6501${CHAIN_ID_DELTACHAIN.toString(
            16
          )}${payload}`,
        },
      };
      vaa.signatures = sign(VAA_SIGNERS, vaa as unknown as VAA<Payload>);
      console.log("uploading", file);
      const msg = client.core.msgStoreCode({
        signer,
        wasm_byte_code: new Uint8Array(contract_bytes),
        vaa: hexToUint8Array(serialiseVAA(vaa as unknown as VAA<Payload>)),
      });
      const result = await client.signAndBroadcast(signer, [msg], {
        ...ZERO_FEE,
        gas: "10000000",
      });
      console.log(result.transactionHash)
      const codeId = Number(
        JSON.parse(result.rawLog)[0]
          .events.find(({ type }) => type === "store_code")
          .attributes.find(({ key }) => key === "code_id").value
      );
      console.log(
        `uploaded ${file}, codeID: ${codeId}, tx: ${result.transactionHash}`
      );

      accum[file] = codeId;
      return accum;
    },
    Object()
  );

  // Instantiate contracts.

  async function instantiate(code_id: number, inst_msg: any, label: string) {
    const instMsgBinary = toBinary(inst_msg);
    const instMsgBytes = fromBase64(instMsgBinary);

    // see /sdk/vaa/governance.go
    const codeIdBuf = Buffer.alloc(8);
    codeIdBuf.writeBigInt64BE(BigInt(code_id));
    const codeIdHash = keccak256(codeIdBuf);
    const codeIdLabelHash = keccak256(
      Buffer.concat([
        Buffer.from(codeIdHash, "hex"),
        Buffer.from(label, "utf8"),
      ])
    );
    const fullHash = keccak256(
      Buffer.concat([Buffer.from(codeIdLabelHash, "hex"), instMsgBytes])
    );

    console.log(fullHash);

    let vaa: VAA<Other> = {
      version: 1,
      phylaxSetIndex: 0,
      signatures: [],
      timestamp: 0,
      nonce: 0,
      emitterChain: GOVERNANCE_CHAIN,
      emitterAddress: GOVERNANCE_EMITTER,
      sequence: BigInt(Math.floor(Math.random() * 100000000)),
      consistencyLevel: 0,
      payload: {
        type: "Other",
        hex: `0000000000000000000000000000000000000000005761736D644D6F64756C6502${CHAIN_ID_DELTACHAIN.toString(
          16
        )}${fullHash}`,
      },
    };
    // TODO: check for number of phylaxs in set and use the corresponding keys
    vaa.signatures = sign(VAA_SIGNERS, vaa as unknown as VAA<Payload>);
    const msg = client.core.msgInstantiateContract({
      signer,
      code_id,
      label,
      msg: instMsgBytes,
      vaa: hexToUint8Array(serialiseVAA(vaa as unknown as VAA<Payload>)),
    });
    const result = await client.signAndBroadcast(signer, [msg], {
      ...ZERO_FEE,
      gas: "10000000",
    });
    console.log("contract instantiation msg: ", msg);
    console.log("contract instantiation result: ", result);
    const addr = JSON.parse(result.rawLog)[0]
      .events.find(({ type }) => type === "instantiate")
      .attributes.find(({ key }) => key === "_contract_address").value;
    console.log(
      `deployed contract ${label}, codeID: ${code_id}, address: ${addr}, txHash: ${result.transactionHash}`
    );

    return addr;
  }

  // Instantiate contracts.
  // NOTE: Only append at the end, the ordering must be deterministic.

  const addresses: {
    [contractName: string]: string;
  } = {};

  const registrations: { [chainName: string]: string } = {
    // keys are only used for logging success/failure
    //solana: String(process.env.REGISTER_SOL_TOKEN_BRIDGE_VAA),
    ethereum: String(process.env.REGISTER_ETH_TOKEN_BRIDGE_VAA),
    bsc: String(process.env.REGISTER_BSC_TOKEN_BRIDGE_VAA),
    //algo: String(process.env.REGISTER_ALGO_TOKEN_BRIDGE_VAA),
    //terra: String(process.env.REGISTER_TERRA_TOKEN_BRIDGE_VAA),
    //near: String(process.env.REGISTER_NEAR_TOKEN_BRIDGE_VAA),
    //terra2: String(process.env.REGISTER_TERRA2_TOKEN_BRIDGE_VAA),
    //aptos: String(process.env.REGISTER_APTOS_TOKEN_BRIDGE_VAA),
    //sui: String(process.env.REGISTER_SUI_TOKEN_BRIDGE_VAA),
    planq: String(process.env.REGISTER_PLANQ_TOKEN_BRIDGE_VAA),
  };

  const instantiateMsg = {};
  addresses["global_accountant.wasm"] = await instantiate(
    codeIds["global_accountant.wasm"],
    instantiateMsg,
    "deltachainAccounting"
  );
  console.log("instantiated accounting: ", addresses["global_accountant.wasm"]);

  const accountingRegistrations = Object.values(registrations).map((r) =>
    Buffer.from(r, "hex").toString("base64")
  );
  const msg = client.wasm.msgExecuteContract({
    sender: signer,
    contract: addresses["global_accountant.wasm"],
    msg: toUtf8(
      JSON.stringify({
        submit_vaas: {
          vaas: accountingRegistrations,
        },
      })
    ),
    funds: [],
  });
  const res = await client.signAndBroadcast(signer, [msg], {
    ...ZERO_FEE,
    gas: "10000000",
  });
  console.log(`sent accounting chain registrations, tx: `, res.transactionHash);

  const deltachainIbcReceiverInstantiateMsg = {};
  addresses["deltachain_ibc_receiver.wasm"] = await instantiate(
    codeIds["deltachain_ibc_receiver.wasm"],
    deltachainIbcReceiverInstantiateMsg,
    "deltachainIbcReceiver"
  );
  console.log(
    "instantiated deltachain ibc receiver contract: ",
    addresses["deltachain_ibc_receiver.wasm"]
  );

  // Generated VAA using
  // `phylaxd template ibc-receiver-update-channel-chain --channel-id channel-0 --chain-id 32 --target-chain-id 7077 > deltachain.prototxt`
  // `phylaxd admin governance-vaa-verify deltachain.prototxt`
  let deltachainIbcReceiverWhitelistVaa: VAA<Other> = {
    version: 1,
    phylaxSetIndex: 0,
    signatures: [],
    timestamp: 0,
    nonce: 0,
    emitterChain: GOVERNANCE_CHAIN,
    emitterAddress: GOVERNANCE_EMITTER,
    sequence: BigInt(Math.floor(Math.random() * 100000000)),
    consistencyLevel: 0,
    payload: {
      type: "Other",
      hex: `0000000000000000000000000000000000000000004962635265636569766572010c20000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006368616e6e656c2d300020`,
    },
  };
  deltachainIbcReceiverWhitelistVaa.signatures = sign(
    VAA_SIGNERS,
    deltachainIbcReceiverWhitelistVaa as unknown as VAA<Payload>
  );
  const deltachainIbcReceiverUpdateWhitelistMsg = {
    submit_update_channel_chain: {
      vaas: [
        Buffer.from(
          serialiseVAA(
            deltachainIbcReceiverWhitelistVaa as unknown as VAA<Payload>
          ),
          "hex"
        ).toString("base64"),
      ],
    },
  };
  const executeMsg = client.wasm.msgExecuteContract({
    sender: signer,
    contract: addresses["deltachain_ibc_receiver.wasm"],
    msg: toUtf8(JSON.stringify(deltachainIbcReceiverUpdateWhitelistMsg)),
    funds: [],
  });
  const updateIbcWhitelistRes = await client.signAndBroadcast(
    signer,
    [executeMsg],
    {
      ...ZERO_FEE,
      gas: "10000000",
    }
  );
  console.log(
    "updated deltachain_ibc_receiver whitelist: ",
    updateIbcWhitelistRes.transactionHash,
    updateIbcWhitelistRes.code
  );
}

try {
  main();
} catch (e: any) {
  if (e?.message) {
    console.error(e.message);
  }
  throw e;
}
