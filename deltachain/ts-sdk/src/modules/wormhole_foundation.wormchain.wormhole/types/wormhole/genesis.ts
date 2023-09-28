//@ts-nocheck
/* eslint-disable */
import { PhylaxSet } from "../deltaswap/phylax_set";
import { Config } from "../deltaswap/config";
import { ReplayProtection } from "../deltaswap/replay_protection";
import { SequenceCounter } from "../deltaswap/sequence_counter";
import { ConsensusPhylaxSetIndex } from "../deltaswap/consensus_phylax_set_index";
import { PhylaxValidator } from "../deltaswap/phylax_validator";
import { Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "deltaswapio.deltachain.deltaswap";

/** GenesisState defines the deltaswap module's genesis state. */
export interface GenesisState {
  phylaxSetList: PhylaxSet[];
  config: Config | undefined;
  replayProtectionList: ReplayProtection[];
  sequenceCounterList: SequenceCounter[];
  consensusPhylaxSetIndex: ConsensusPhylaxSetIndex | undefined;
  /** this line is used by starport scaffolding # genesis/proto/state */
  phylaxValidatorList: PhylaxValidator[];
}

const baseGenesisState: object = {};

export const GenesisState = {
  encode(message: GenesisState, writer: Writer = Writer.create()): Writer {
    for (const v of message.phylaxSetList) {
      PhylaxSet.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.config !== undefined) {
      Config.encode(message.config, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.replayProtectionList) {
      ReplayProtection.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.sequenceCounterList) {
      SequenceCounter.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    if (message.consensusPhylaxSetIndex !== undefined) {
      ConsensusPhylaxSetIndex.encode(
        message.consensusPhylaxSetIndex,
        writer.uint32(42).fork()
      ).ldelim();
    }
    for (const v of message.phylaxValidatorList) {
      PhylaxValidator.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseGenesisState } as GenesisState;
    message.phylaxSetList = [];
    message.replayProtectionList = [];
    message.sequenceCounterList = [];
    message.phylaxValidatorList = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.phylaxSetList.push(
            PhylaxSet.decode(reader, reader.uint32())
          );
          break;
        case 2:
          message.config = Config.decode(reader, reader.uint32());
          break;
        case 3:
          message.replayProtectionList.push(
            ReplayProtection.decode(reader, reader.uint32())
          );
          break;
        case 4:
          message.sequenceCounterList.push(
            SequenceCounter.decode(reader, reader.uint32())
          );
          break;
        case 5:
          message.consensusPhylaxSetIndex = ConsensusPhylaxSetIndex.decode(
            reader,
            reader.uint32()
          );
          break;
        case 6:
          message.phylaxValidatorList.push(
            PhylaxValidator.decode(reader, reader.uint32())
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GenesisState {
    const message = { ...baseGenesisState } as GenesisState;
    message.phylaxSetList = [];
    message.replayProtectionList = [];
    message.sequenceCounterList = [];
    message.phylaxValidatorList = [];
    if (
      object.phylaxSetList !== undefined &&
      object.phylaxSetList !== null
    ) {
      for (const e of object.phylaxSetList) {
        message.phylaxSetList.push(PhylaxSet.fromJSON(e));
      }
    }
    if (object.config !== undefined && object.config !== null) {
      message.config = Config.fromJSON(object.config);
    } else {
      message.config = undefined;
    }
    if (
      object.replayProtectionList !== undefined &&
      object.replayProtectionList !== null
    ) {
      for (const e of object.replayProtectionList) {
        message.replayProtectionList.push(ReplayProtection.fromJSON(e));
      }
    }
    if (
      object.sequenceCounterList !== undefined &&
      object.sequenceCounterList !== null
    ) {
      for (const e of object.sequenceCounterList) {
        message.sequenceCounterList.push(SequenceCounter.fromJSON(e));
      }
    }
    if (
      object.consensusPhylaxSetIndex !== undefined &&
      object.consensusPhylaxSetIndex !== null
    ) {
      message.consensusPhylaxSetIndex = ConsensusPhylaxSetIndex.fromJSON(
        object.consensusPhylaxSetIndex
      );
    } else {
      message.consensusPhylaxSetIndex = undefined;
    }
    if (
      object.phylaxValidatorList !== undefined &&
      object.phylaxValidatorList !== null
    ) {
      for (const e of object.phylaxValidatorList) {
        message.phylaxValidatorList.push(PhylaxValidator.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    if (message.phylaxSetList) {
      obj.phylaxSetList = message.phylaxSetList.map((e) =>
        e ? PhylaxSet.toJSON(e) : undefined
      );
    } else {
      obj.phylaxSetList = [];
    }
    message.config !== undefined &&
      (obj.config = message.config ? Config.toJSON(message.config) : undefined);
    if (message.replayProtectionList) {
      obj.replayProtectionList = message.replayProtectionList.map((e) =>
        e ? ReplayProtection.toJSON(e) : undefined
      );
    } else {
      obj.replayProtectionList = [];
    }
    if (message.sequenceCounterList) {
      obj.sequenceCounterList = message.sequenceCounterList.map((e) =>
        e ? SequenceCounter.toJSON(e) : undefined
      );
    } else {
      obj.sequenceCounterList = [];
    }
    message.consensusPhylaxSetIndex !== undefined &&
      (obj.consensusPhylaxSetIndex = message.consensusPhylaxSetIndex
        ? ConsensusPhylaxSetIndex.toJSON(message.consensusPhylaxSetIndex)
        : undefined);
    if (message.phylaxValidatorList) {
      obj.phylaxValidatorList = message.phylaxValidatorList.map((e) =>
        e ? PhylaxValidator.toJSON(e) : undefined
      );
    } else {
      obj.phylaxValidatorList = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<GenesisState>): GenesisState {
    const message = { ...baseGenesisState } as GenesisState;
    message.phylaxSetList = [];
    message.replayProtectionList = [];
    message.sequenceCounterList = [];
    message.phylaxValidatorList = [];
    if (
      object.phylaxSetList !== undefined &&
      object.phylaxSetList !== null
    ) {
      for (const e of object.phylaxSetList) {
        message.phylaxSetList.push(PhylaxSet.fromPartial(e));
      }
    }
    if (object.config !== undefined && object.config !== null) {
      message.config = Config.fromPartial(object.config);
    } else {
      message.config = undefined;
    }
    if (
      object.replayProtectionList !== undefined &&
      object.replayProtectionList !== null
    ) {
      for (const e of object.replayProtectionList) {
        message.replayProtectionList.push(ReplayProtection.fromPartial(e));
      }
    }
    if (
      object.sequenceCounterList !== undefined &&
      object.sequenceCounterList !== null
    ) {
      for (const e of object.sequenceCounterList) {
        message.sequenceCounterList.push(SequenceCounter.fromPartial(e));
      }
    }
    if (
      object.consensusPhylaxSetIndex !== undefined &&
      object.consensusPhylaxSetIndex !== null
    ) {
      message.consensusPhylaxSetIndex = ConsensusPhylaxSetIndex.fromPartial(
        object.consensusPhylaxSetIndex
      );
    } else {
      message.consensusPhylaxSetIndex = undefined;
    }
    if (
      object.phylaxValidatorList !== undefined &&
      object.phylaxValidatorList !== null
    ) {
      for (const e of object.phylaxValidatorList) {
        message.phylaxValidatorList.push(PhylaxValidator.fromPartial(e));
      }
    }
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | undefined;
export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;
