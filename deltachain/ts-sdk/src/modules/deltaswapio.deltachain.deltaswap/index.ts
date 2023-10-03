//@ts-nocheck
// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

import { StdFee } from "@cosmjs/launchpad";
import { SigningStargateClient } from "@cosmjs/stargate";
import { Registry, OfflineSigner, EncodeObject, DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
import { Api } from "./rest";
import { MsgExecuteGovernanceVAA } from "./types/deltaswap/tx";
import { MsgStoreCode } from "./types/deltaswap/tx";
import { MsgCreateAllowlistEntryRequest } from "./types/deltaswap/tx";
import { MsgDeleteAllowlistEntryRequest } from "./types/deltaswap/tx";
import { MsgInstantiateContract } from "./types/deltaswap/tx";
import { MsgExecuteGatewayGovernanceVaa } from "./types/deltaswap/tx";
import { MsgDeleteWasmInstantiateAllowlist } from "./types/deltaswap/tx";
import { MsgMigrateContract } from "./types/deltaswap/tx";
import { MsgAddWasmInstantiateAllowlist } from "./types/deltaswap/tx";
import { MsgRegisterAccountAsPhylax } from "./types/deltaswap/tx";


const types = [
  ["/deltaswapio.deltachain.deltaswap.MsgExecuteGovernanceVAA", MsgExecuteGovernanceVAA],
  ["/deltaswapio.deltachain.deltaswap.MsgStoreCode", MsgStoreCode],
  ["/deltaswapio.deltachain.deltaswap.MsgCreateAllowlistEntryRequest", MsgCreateAllowlistEntryRequest],
  ["/deltaswapio.deltachain.deltaswap.MsgDeleteAllowlistEntryRequest", MsgDeleteAllowlistEntryRequest],
  ["/deltaswapio.deltachain.deltaswap.MsgInstantiateContract", MsgInstantiateContract],
  ["/deltaswapio.deltachain.deltaswap.MsgExecuteGatewayGovernanceVaa", MsgExecuteGatewayGovernanceVaa],
  ["/deltaswapio.deltachain.deltaswap.MsgDeleteWasmInstantiateAllowlist", MsgDeleteWasmInstantiateAllowlist],
  ["/deltaswapio.deltachain.deltaswap.MsgMigrateContract", MsgMigrateContract],
  ["/deltaswapio.deltachain.deltaswap.MsgAddWasmInstantiateAllowlist", MsgAddWasmInstantiateAllowlist],
  ["/deltaswapio.deltachain.deltaswap.MsgRegisterAccountAsPhylax", MsgRegisterAccountAsPhylax],
  
];
export const MissingWalletError = new Error("wallet is required");

export const registry = new Registry(<any>types);

const defaultFee = {
  amount: [],
  gas: "200000",
};

interface TxClientOptions {
  addr: string
}

interface SignAndBroadcastOptions {
  fee: StdFee,
  memo?: string
}

const txClient = async (wallet: OfflineSigner, { addr: addr }: TxClientOptions = { addr: "http://localhost:26657" }) => {
  if (!wallet) throw MissingWalletError;
  let client;
  if (addr) {
    client = await SigningStargateClient.connectWithSigner(addr, wallet, { registry });
  }else{
    client = await SigningStargateClient.offline( wallet, { registry });
  }
  const { address } = (await wallet.getAccounts())[0];

  return {
    signAndBroadcast: (msgs: EncodeObject[], { fee, memo }: SignAndBroadcastOptions = {fee: defaultFee, memo: ""}) => client.signAndBroadcast(address, msgs, fee,memo),
    msgExecuteGovernanceVAA: (data: MsgExecuteGovernanceVAA): EncodeObject => ({ typeUrl: "/deltaswapio.deltachain.deltaswap.MsgExecuteGovernanceVAA", value: MsgExecuteGovernanceVAA.fromPartial( data ) }),
    msgStoreCode: (data: MsgStoreCode): EncodeObject => ({ typeUrl: "/deltaswapio.deltachain.deltaswap.MsgStoreCode", value: MsgStoreCode.fromPartial( data ) }),
    msgCreateAllowlistEntryRequest: (data: MsgCreateAllowlistEntryRequest): EncodeObject => ({ typeUrl: "/deltaswapio.deltachain.deltaswap.MsgCreateAllowlistEntryRequest", value: MsgCreateAllowlistEntryRequest.fromPartial( data ) }),
    msgDeleteAllowlistEntryRequest: (data: MsgDeleteAllowlistEntryRequest): EncodeObject => ({ typeUrl: "/deltaswapio.deltachain.deltaswap.MsgDeleteAllowlistEntryRequest", value: MsgDeleteAllowlistEntryRequest.fromPartial( data ) }),
    msgInstantiateContract: (data: MsgInstantiateContract): EncodeObject => ({ typeUrl: "/deltaswapio.deltachain.deltaswap.MsgInstantiateContract", value: MsgInstantiateContract.fromPartial( data ) }),
    msgExecuteGatewayGovernanceVaa: (data: MsgExecuteGatewayGovernanceVaa): EncodeObject => ({ typeUrl: "/deltaswapio.deltachain.deltaswap.MsgExecuteGatewayGovernanceVaa", value: MsgExecuteGatewayGovernanceVaa.fromPartial( data ) }),
    msgDeleteWasmInstantiateAllowlist: (data: MsgDeleteWasmInstantiateAllowlist): EncodeObject => ({ typeUrl: "/deltaswapio.deltachain.deltaswap.MsgDeleteWasmInstantiateAllowlist", value: MsgDeleteWasmInstantiateAllowlist.fromPartial( data ) }),
    msgMigrateContract: (data: MsgMigrateContract): EncodeObject => ({ typeUrl: "/deltaswapio.deltachain.deltaswap.MsgMigrateContract", value: MsgMigrateContract.fromPartial( data ) }),
    msgAddWasmInstantiateAllowlist: (data: MsgAddWasmInstantiateAllowlist): EncodeObject => ({ typeUrl: "/deltaswapio.deltachain.deltaswap.MsgAddWasmInstantiateAllowlist", value: MsgAddWasmInstantiateAllowlist.fromPartial( data ) }),
    msgRegisterAccountAsPhylax: (data: MsgRegisterAccountAsPhylax): EncodeObject => ({ typeUrl: "/deltaswapio.deltachain.deltaswap.MsgRegisterAccountAsPhylax", value: MsgRegisterAccountAsPhylax.fromPartial( data ) }),
    
  };
};

interface QueryClientOptions {
  addr: string
}

const queryClient = async ({ addr: addr }: QueryClientOptions = { addr: "http://localhost:1317" }) => {
  return new Api({ baseUrl: addr });
};

export {
  txClient,
  queryClient,
};
