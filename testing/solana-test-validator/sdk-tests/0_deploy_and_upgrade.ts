import { expect } from "chai";
import * as web3 from "@solana/web3.js";
import fs from "fs";
import { MockPhylaxs, GovernanceEmitter } from "../../../sdk/js/src/mock";
import {
  getPhylaxSet,
  getDeltaswapBridgeData,
  createInitializeInstruction as createDeltaswapInitializeInstruction,
  deriveUpgradeAuthorityKey,
  createUpgradeContractInstruction as createDeltaswapUpgradeContractInstruction,
} from "../../../sdk/js/src/solana/wormhole";
import {
  createInitializeInstruction as createTokenBridgeInitializeInstruction,
  createUpgradeContractInstruction as createTokenBridgeUpgradeContractInstruction,
  getTokenBridgeConfig,
} from "../../../sdk/js/src/solana/tokenBridge";
import {
  createInitializeInstruction as createNftBridgeInitializeInstruction,
  createUpgradeContractInstruction as createNftBridgeUpgradeContractInstruction,
} from "../../../sdk/js/src/solana/nftBridge";
import { postVaa } from "../../../sdk/js/src/solana/sendAndConfirmPostVaa";
import { NodeWallet } from "../../../sdk/js/src/solana/utils";

import {
  CORE_BRIDGE_ADDRESS,
  GOVERNANCE_EMITTER_ADDRESS,
  PHYLAX_KEYS,
  PHYLAX_SET_INDEX,
  LOCALHOST,
  NFT_BRIDGE_ADDRESS,
  TOKEN_BRIDGE_ADDRESS,
} from "./helpers/consts";
import {
  deployProgram,
  execSolanaWriteBufferAndSetBufferAuthority,
} from "./helpers/utils";
import { getNftBridgeConfig } from "../../../sdk/js/src/solana/nftBridge";

describe("Deploy and Upgrade Programs", () => {
  const connection = new web3.Connection(LOCALHOST, "processed");

  const payerPath = `${__dirname}/keys/pFCBP4bhqdSsrWUVTgqhPsLrfEdChBK17vgFM7TxjxQ.json`;
  const wallet = NodeWallet.fromSecretKey(
    Uint8Array.from(
      JSON.parse(
        fs.readFileSync(payerPath, {
          encoding: "utf8",
        })
      )
    )
  );

  // for signing wormhole messages
  const phylaxs = new MockPhylaxs(PHYLAX_SET_INDEX, PHYLAX_KEYS);

  // for generating governance wormhole messages
  const governance = new GovernanceEmitter(
    GOVERNANCE_EMITTER_ADDRESS.toBuffer().toString("hex")
  );

  const localVariables: any = {};

  before("Airdrop SOL", async () => {
    // wallet
    await connection
      .requestAirdrop(wallet.key(), 1000 * web3.LAMPORTS_PER_SOL)
      .then(async (signature) => connection.confirmTransaction(signature));
  });

  describe("Deltaswap (Core Bridge)", () => {
    it("Deploy and Initialize", async () => {
      const artifactPath = `${__dirname}/../artifacts-main/bridge.so`;
      const programIdPath = `${__dirname}/keys/${CORE_BRIDGE_ADDRESS}.json`;
      const upgradeAuthority = deriveUpgradeAuthorityKey(CORE_BRIDGE_ADDRESS);

      deployProgram(
        payerPath,
        artifactPath,
        programIdPath,
        CORE_BRIDGE_ADDRESS,
        upgradeAuthority
      );

      // initialize
      const phylaxSetExpirationTime = 86400;
      const fee = 100n;
      const initialPhylaxs = phylaxs.getPublicKeys().slice(0, 1);

      const initializeTx = await web3.sendAndConfirmTransaction(
        connection,
        new web3.Transaction().add(
          createDeltaswapInitializeInstruction(
            CORE_BRIDGE_ADDRESS,
            wallet.key(),
            phylaxSetExpirationTime,
            fee,
            initialPhylaxs
          )
        ),
        [wallet.signer()]
      );
      // console.log(`initializeTx: ${initializeTx}`);

      // verify data
      const info = await getDeltaswapBridgeData(connection, CORE_BRIDGE_ADDRESS);
      expect(info.phylaxSetIndex).to.equal(PHYLAX_SET_INDEX);
      expect(info.config.phylaxSetExpirationTime).to.equal(
        phylaxSetExpirationTime
      );
      expect(info.config.fee).to.equal(fee);
      const phylaxSet = await getPhylaxSet(
        connection,
        CORE_BRIDGE_ADDRESS,
        0
      );
      expect(phylaxSet.index).to.equal(PHYLAX_SET_INDEX);
      expect(phylaxSet.expirationTime).to.equal(0);
      expect(phylaxSet.keys).has.length(1);
      expect(
        Buffer.compare(initialPhylaxs[0], phylaxSet.keys.at(0)!)
      ).to.equal(0);
    });

    it("Upgrade Contract", async () => {
      // first upload BPF of implementation
      const artifactPath = `${__dirname}/../artifacts/bridge.so`;
      const upgradeAuthority = deriveUpgradeAuthorityKey(CORE_BRIDGE_ADDRESS);

      const implementation = execSolanaWriteBufferAndSetBufferAuthority(
        payerPath,
        artifactPath,
        upgradeAuthority
      );

      // now pass implementation through governance
      const timestamp = 1;
      const chain = 1;
      const message = governance.publishDeltaswapUpgradeContract(
        timestamp,
        chain,
        implementation.toString()
      );
      const signedVaa = phylaxs.addSignatures(message, [0]);
      // console.log(`signedVaa: ${signedVaa.toString("base64")}`);

      const txSignatures = await postVaa(
        connection,
        wallet.signTransaction,
        CORE_BRIDGE_ADDRESS,
        wallet.key(),
        signedVaa
      ).then((results) => results.map((result) => result.signature));
      const postTx = txSignatures.pop()!;
      for (const verifyTx of txSignatures) {
        // console.log(`verifySignatures:  ${verifyTx}`);
      }
      // console.log(`postVaa:           ${postTx}`);

      const upgradeContractIx = createDeltaswapUpgradeContractInstruction(
        CORE_BRIDGE_ADDRESS,
        wallet.key(),
        signedVaa
      );

      const upgradeContractTx = await web3.sendAndConfirmTransaction(
        connection,
        new web3.Transaction().add(upgradeContractIx),
        [wallet.signer()]
      );
      // console.log(`upgradeContractTx: ${upgradeContractTx}`);

      // not sure what to verify. if there are no errors,
      // transaction was successful and the upgrade test passes
    });
  });

  describe("Token Bridge", () => {
    it("Deploy and Initialize", async () => {
      const artifactPath = `${__dirname}/../artifacts-main/token_bridge.so`;
      const programIdPath = `${__dirname}/keys/${TOKEN_BRIDGE_ADDRESS}.json`;
      const upgradeAuthority = deriveUpgradeAuthorityKey(TOKEN_BRIDGE_ADDRESS);

      deployProgram(
        payerPath,
        artifactPath,
        programIdPath,
        TOKEN_BRIDGE_ADDRESS,
        upgradeAuthority
      );

      // we will initialize using CORE_BRIDGE_ADDRESS instead of
      // UPGRADEABLE_CORE_BRIDGE_ADDRESS because the Deltaswap owner is only
      // valid for CORE_BRIDGE_ADDRESS with the bpf we deployed
      //
      // AccountOwner::Other(Pubkey::from_str(env!("BRIDGE_ADDRESS")).unwrap())
      const initializeTx = await web3.sendAndConfirmTransaction(
        connection,
        new web3.Transaction().add(
          createTokenBridgeInitializeInstruction(
            TOKEN_BRIDGE_ADDRESS,
            wallet.key(),
            CORE_BRIDGE_ADDRESS
          )
        ),
        [wallet.signer()]
      );
      // console.log(`initializeTx: ${initializeTx}`);

      // verify data
      const config = await getTokenBridgeConfig(
        connection,
        TOKEN_BRIDGE_ADDRESS
      );
      expect(config.wormhole.equals(CORE_BRIDGE_ADDRESS)).to.be.true;
    });

    it("Upgrade Contract", async () => {
      // first upload BPF of implementation
      const artifactPath = `${__dirname}/../artifacts/token_bridge.so`;
      const upgradeAuthority = deriveUpgradeAuthorityKey(TOKEN_BRIDGE_ADDRESS);

      const implementation = execSolanaWriteBufferAndSetBufferAuthority(
        payerPath,
        artifactPath,
        upgradeAuthority
      );

      // now pass implementation through governance
      const timestamp = 2;
      const chain = 1;
      const message = governance.publishTokenBridgeUpgradeContract(
        timestamp,
        chain,
        implementation.toString()
      );
      const signedVaa = phylaxs.addSignatures(message, [0]);
      // console.log(`signedVaa: ${signedVaa.toString("base64")}`);

      const txSignatures = await postVaa(
        connection,
        wallet.signTransaction,
        CORE_BRIDGE_ADDRESS,
        wallet.key(),
        signedVaa
      ).then((results) => results.map((result) => result.signature));
      const postTx = txSignatures.pop()!;
      for (const verifyTx of txSignatures) {
        // console.log(`verifySignatures:  ${verifyTx}`);
      }
      // console.log(`postVaa:           ${postTx}`);

      const upgradeContractIx = createTokenBridgeUpgradeContractInstruction(
        TOKEN_BRIDGE_ADDRESS,
        CORE_BRIDGE_ADDRESS,
        wallet.key(),
        signedVaa
      );

      const upgradeContractTx = await web3.sendAndConfirmTransaction(
        connection,
        new web3.Transaction().add(upgradeContractIx),
        [wallet.signer()]
      );
      // console.log(`upgradeContractTx: ${upgradeContractTx}`);

      // not sure what to verify. if there are no errors,
      // transaction was successful and the upgrade test passes
    });
  });

  describe("NFT Bridge", () => {
    it("Deploy and Initialize", async () => {
      const artifactPath = `${__dirname}/../artifacts-main/nft_bridge.so`;
      const programIdPath = `${__dirname}/keys/${NFT_BRIDGE_ADDRESS}.json`;
      const upgradeAuthority = deriveUpgradeAuthorityKey(NFT_BRIDGE_ADDRESS);

      deployProgram(
        payerPath,
        artifactPath,
        programIdPath,
        NFT_BRIDGE_ADDRESS,
        upgradeAuthority
      );

      // we will initialize using CORE_BRIDGE_ADDRESS instead of
      // UPGRADEABLE_CORE_BRIDGE_ADDRESS because the Deltaswap owner is only
      // valid for CORE_BRIDGE_ADDRESS with the bpf we deployed
      //
      // AccountOwner::Other(Pubkey::from_str(env!("BRIDGE_ADDRESS")).unwrap())
      const initializeTx = await web3.sendAndConfirmTransaction(
        connection,
        new web3.Transaction().add(
          createNftBridgeInitializeInstruction(
            NFT_BRIDGE_ADDRESS,
            wallet.key(),
            CORE_BRIDGE_ADDRESS
          )
        ),
        [wallet.signer()]
      );
      // console.log(`initializeTx: ${initializeTx}`);

      // verify data
      const config = await getNftBridgeConfig(connection, NFT_BRIDGE_ADDRESS);
      expect(config.wormhole.equals(CORE_BRIDGE_ADDRESS)).to.be.true;
    });

    it("Upgrade Contract", async () => {
      // first upload BPF of implementation
      const artifactPath = `${__dirname}/../artifacts/nft_bridge.so`;
      const upgradeAuthority = deriveUpgradeAuthorityKey(NFT_BRIDGE_ADDRESS);

      const implementation = execSolanaWriteBufferAndSetBufferAuthority(
        payerPath,
        artifactPath,
        upgradeAuthority
      );

      // now pass implementation through governance
      const timestamp = 3;
      const chain = 1;
      const message = governance.publishNftBridgeUpgradeContract(
        timestamp,
        chain,
        implementation.toString()
      );
      const signedVaa = phylaxs.addSignatures(message, [0]);
      // console.log(`signedVaa: ${signedVaa.toString("base64")}`);

      const txSignatures = await postVaa(
        connection,
        wallet.signTransaction,
        CORE_BRIDGE_ADDRESS,
        wallet.key(),
        signedVaa
      ).then((results) => results.map((result) => result.signature));
      const postTx = txSignatures.pop()!;
      for (const verifyTx of txSignatures) {
        // console.log(`verifySignatures:  ${verifyTx}`);
      }
      // console.log(`postVaa:           ${postTx}`);

      const upgradeContractIx = createNftBridgeUpgradeContractInstruction(
        NFT_BRIDGE_ADDRESS,
        CORE_BRIDGE_ADDRESS,
        wallet.key(),
        signedVaa
      );

      const upgradeContractTx = await web3.sendAndConfirmTransaction(
        connection,
        new web3.Transaction().add(upgradeContractIx),
        [wallet.signer()]
      );
      // console.log(`upgradeContractTx: ${upgradeContractTx}`);

      // not sure what to verify. if there are no errors,
      // transaction was successful and the upgrade test passes
    });
  });
});
