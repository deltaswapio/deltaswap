//@ts-nocheck
/* eslint-disable */
import { Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "deltaswapio.deltachain.deltaswap";

export interface ConsensusPhylaxSetIndex {
  index: number;
}

const baseConsensusPhylaxSetIndex: object = { index: 0 };

export const ConsensusPhylaxSetIndex = {
  encode(
    message: ConsensusPhylaxSetIndex,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.index !== 0) {
      writer.uint32(8).uint32(message.index);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): ConsensusPhylaxSetIndex {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseConsensusPhylaxSetIndex,
    } as ConsensusPhylaxSetIndex;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.index = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ConsensusPhylaxSetIndex {
    const message = {
      ...baseConsensusPhylaxSetIndex,
    } as ConsensusPhylaxSetIndex;
    if (object.index !== undefined && object.index !== null) {
      message.index = Number(object.index);
    } else {
      message.index = 0;
    }
    return message;
  },

  toJSON(message: ConsensusPhylaxSetIndex): unknown {
    const obj: any = {};
    message.index !== undefined && (obj.index = message.index);
    return obj;
  },

  fromPartial(
    object: DeepPartial<ConsensusPhylaxSetIndex>
  ): ConsensusPhylaxSetIndex {
    const message = {
      ...baseConsensusPhylaxSetIndex,
    } as ConsensusPhylaxSetIndex;
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    } else {
      message.index = 0;
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
