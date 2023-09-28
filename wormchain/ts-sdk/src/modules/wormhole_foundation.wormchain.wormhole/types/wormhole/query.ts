//@ts-nocheck
/* eslint-disable */
import { Reader, Writer } from "protobufjs/minimal";
import { PhylaxSet } from "../wormhole/guardian_set";
import {
  PageRequest,
  PageResponse,
} from "../cosmos/base/query/v1beta1/pagination";
import { Config } from "../wormhole/config";
import { ReplayProtection } from "../wormhole/replay_protection";
import { SequenceCounter } from "../wormhole/sequence_counter";
import { ConsensusPhylaxSetIndex } from "../wormhole/consensus_guardian_set_index";
import { PhylaxValidator } from "../wormhole/guardian_validator";

export const protobufPackage = "wormhole_foundation.deltachain.wormhole";

export interface QueryGetPhylaxSetRequest {
  index: number;
}

export interface QueryGetPhylaxSetResponse {
  PhylaxSet: PhylaxSet | undefined;
}

export interface QueryAllPhylaxSetRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllPhylaxSetResponse {
  PhylaxSet: PhylaxSet[];
  pagination: PageResponse | undefined;
}

export interface QueryGetConfigRequest {}

export interface QueryGetConfigResponse {
  Config: Config | undefined;
}

export interface QueryGetReplayProtectionRequest {
  index: string;
}

export interface QueryGetReplayProtectionResponse {
  replayProtection: ReplayProtection | undefined;
}

export interface QueryAllReplayProtectionRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllReplayProtectionResponse {
  replayProtection: ReplayProtection[];
  pagination: PageResponse | undefined;
}

export interface QueryGetSequenceCounterRequest {
  index: string;
}

export interface QueryGetSequenceCounterResponse {
  sequenceCounter: SequenceCounter | undefined;
}

export interface QueryAllSequenceCounterRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllSequenceCounterResponse {
  sequenceCounter: SequenceCounter[];
  pagination: PageResponse | undefined;
}

export interface QueryGetConsensusPhylaxSetIndexRequest {}

export interface QueryGetConsensusPhylaxSetIndexResponse {
  ConsensusPhylaxSetIndex: ConsensusPhylaxSetIndex | undefined;
}

export interface QueryGetPhylaxValidatorRequest {
  guardianKey: Uint8Array;
}

export interface QueryGetPhylaxValidatorResponse {
  guardianValidator: PhylaxValidator | undefined;
}

export interface QueryAllPhylaxValidatorRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllPhylaxValidatorResponse {
  guardianValidator: PhylaxValidator[];
  pagination: PageResponse | undefined;
}

export interface QueryLatestPhylaxSetIndexRequest {}

export interface QueryLatestPhylaxSetIndexResponse {
  latestPhylaxSetIndex: number;
}

const baseQueryGetPhylaxSetRequest: object = { index: 0 };

export const QueryGetPhylaxSetRequest = {
  encode(
    message: QueryGetPhylaxSetRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.index !== 0) {
      writer.uint32(8).uint32(message.index);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryGetPhylaxSetRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryGetPhylaxSetRequest,
    } as QueryGetPhylaxSetRequest;
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

  fromJSON(object: any): QueryGetPhylaxSetRequest {
    const message = {
      ...baseQueryGetPhylaxSetRequest,
    } as QueryGetPhylaxSetRequest;
    if (object.index !== undefined && object.index !== null) {
      message.index = Number(object.index);
    } else {
      message.index = 0;
    }
    return message;
  },

  toJSON(message: QueryGetPhylaxSetRequest): unknown {
    const obj: any = {};
    message.index !== undefined && (obj.index = message.index);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetPhylaxSetRequest>
  ): QueryGetPhylaxSetRequest {
    const message = {
      ...baseQueryGetPhylaxSetRequest,
    } as QueryGetPhylaxSetRequest;
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    } else {
      message.index = 0;
    }
    return message;
  },
};

const baseQueryGetPhylaxSetResponse: object = {};

export const QueryGetPhylaxSetResponse = {
  encode(
    message: QueryGetPhylaxSetResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.PhylaxSet !== undefined) {
      PhylaxSet.encode(
        message.PhylaxSet,
        writer.uint32(10).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryGetPhylaxSetResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryGetPhylaxSetResponse,
    } as QueryGetPhylaxSetResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.PhylaxSet = PhylaxSet.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetPhylaxSetResponse {
    const message = {
      ...baseQueryGetPhylaxSetResponse,
    } as QueryGetPhylaxSetResponse;
    if (object.PhylaxSet !== undefined && object.PhylaxSet !== null) {
      message.PhylaxSet = PhylaxSet.fromJSON(object.PhylaxSet);
    } else {
      message.PhylaxSet = undefined;
    }
    return message;
  },

  toJSON(message: QueryGetPhylaxSetResponse): unknown {
    const obj: any = {};
    message.PhylaxSet !== undefined &&
      (obj.PhylaxSet = message.PhylaxSet
        ? PhylaxSet.toJSON(message.PhylaxSet)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetPhylaxSetResponse>
  ): QueryGetPhylaxSetResponse {
    const message = {
      ...baseQueryGetPhylaxSetResponse,
    } as QueryGetPhylaxSetResponse;
    if (object.PhylaxSet !== undefined && object.PhylaxSet !== null) {
      message.PhylaxSet = PhylaxSet.fromPartial(object.PhylaxSet);
    } else {
      message.PhylaxSet = undefined;
    }
    return message;
  },
};

const baseQueryAllPhylaxSetRequest: object = {};

export const QueryAllPhylaxSetRequest = {
  encode(
    message: QueryAllPhylaxSetRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryAllPhylaxSetRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllPhylaxSetRequest,
    } as QueryAllPhylaxSetRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllPhylaxSetRequest {
    const message = {
      ...baseQueryAllPhylaxSetRequest,
    } as QueryAllPhylaxSetRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllPhylaxSetRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageRequest.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllPhylaxSetRequest>
  ): QueryAllPhylaxSetRequest {
    const message = {
      ...baseQueryAllPhylaxSetRequest,
    } as QueryAllPhylaxSetRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryAllPhylaxSetResponse: object = {};

export const QueryAllPhylaxSetResponse = {
  encode(
    message: QueryAllPhylaxSetResponse,
    writer: Writer = Writer.create()
  ): Writer {
    for (const v of message.PhylaxSet) {
      PhylaxSet.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(
        message.pagination,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryAllPhylaxSetResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllPhylaxSetResponse,
    } as QueryAllPhylaxSetResponse;
    message.PhylaxSet = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.PhylaxSet.push(PhylaxSet.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllPhylaxSetResponse {
    const message = {
      ...baseQueryAllPhylaxSetResponse,
    } as QueryAllPhylaxSetResponse;
    message.PhylaxSet = [];
    if (object.PhylaxSet !== undefined && object.PhylaxSet !== null) {
      for (const e of object.PhylaxSet) {
        message.PhylaxSet.push(PhylaxSet.fromJSON(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllPhylaxSetResponse): unknown {
    const obj: any = {};
    if (message.PhylaxSet) {
      obj.PhylaxSet = message.PhylaxSet.map((e) =>
        e ? PhylaxSet.toJSON(e) : undefined
      );
    } else {
      obj.PhylaxSet = [];
    }
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageResponse.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllPhylaxSetResponse>
  ): QueryAllPhylaxSetResponse {
    const message = {
      ...baseQueryAllPhylaxSetResponse,
    } as QueryAllPhylaxSetResponse;
    message.PhylaxSet = [];
    if (object.PhylaxSet !== undefined && object.PhylaxSet !== null) {
      for (const e of object.PhylaxSet) {
        message.PhylaxSet.push(PhylaxSet.fromPartial(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryGetConfigRequest: object = {};

export const QueryGetConfigRequest = {
  encode(_: QueryGetConfigRequest, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryGetConfigRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryGetConfigRequest } as QueryGetConfigRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): QueryGetConfigRequest {
    const message = { ...baseQueryGetConfigRequest } as QueryGetConfigRequest;
    return message;
  },

  toJSON(_: QueryGetConfigRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<QueryGetConfigRequest>): QueryGetConfigRequest {
    const message = { ...baseQueryGetConfigRequest } as QueryGetConfigRequest;
    return message;
  },
};

const baseQueryGetConfigResponse: object = {};

export const QueryGetConfigResponse = {
  encode(
    message: QueryGetConfigResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.Config !== undefined) {
      Config.encode(message.Config, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryGetConfigResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryGetConfigResponse } as QueryGetConfigResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.Config = Config.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetConfigResponse {
    const message = { ...baseQueryGetConfigResponse } as QueryGetConfigResponse;
    if (object.Config !== undefined && object.Config !== null) {
      message.Config = Config.fromJSON(object.Config);
    } else {
      message.Config = undefined;
    }
    return message;
  },

  toJSON(message: QueryGetConfigResponse): unknown {
    const obj: any = {};
    message.Config !== undefined &&
      (obj.Config = message.Config ? Config.toJSON(message.Config) : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetConfigResponse>
  ): QueryGetConfigResponse {
    const message = { ...baseQueryGetConfigResponse } as QueryGetConfigResponse;
    if (object.Config !== undefined && object.Config !== null) {
      message.Config = Config.fromPartial(object.Config);
    } else {
      message.Config = undefined;
    }
    return message;
  },
};

const baseQueryGetReplayProtectionRequest: object = { index: "" };

export const QueryGetReplayProtectionRequest = {
  encode(
    message: QueryGetReplayProtectionRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.index !== "") {
      writer.uint32(10).string(message.index);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryGetReplayProtectionRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryGetReplayProtectionRequest,
    } as QueryGetReplayProtectionRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.index = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetReplayProtectionRequest {
    const message = {
      ...baseQueryGetReplayProtectionRequest,
    } as QueryGetReplayProtectionRequest;
    if (object.index !== undefined && object.index !== null) {
      message.index = String(object.index);
    } else {
      message.index = "";
    }
    return message;
  },

  toJSON(message: QueryGetReplayProtectionRequest): unknown {
    const obj: any = {};
    message.index !== undefined && (obj.index = message.index);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetReplayProtectionRequest>
  ): QueryGetReplayProtectionRequest {
    const message = {
      ...baseQueryGetReplayProtectionRequest,
    } as QueryGetReplayProtectionRequest;
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    } else {
      message.index = "";
    }
    return message;
  },
};

const baseQueryGetReplayProtectionResponse: object = {};

export const QueryGetReplayProtectionResponse = {
  encode(
    message: QueryGetReplayProtectionResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.replayProtection !== undefined) {
      ReplayProtection.encode(
        message.replayProtection,
        writer.uint32(10).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryGetReplayProtectionResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryGetReplayProtectionResponse,
    } as QueryGetReplayProtectionResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.replayProtection = ReplayProtection.decode(
            reader,
            reader.uint32()
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetReplayProtectionResponse {
    const message = {
      ...baseQueryGetReplayProtectionResponse,
    } as QueryGetReplayProtectionResponse;
    if (
      object.replayProtection !== undefined &&
      object.replayProtection !== null
    ) {
      message.replayProtection = ReplayProtection.fromJSON(
        object.replayProtection
      );
    } else {
      message.replayProtection = undefined;
    }
    return message;
  },

  toJSON(message: QueryGetReplayProtectionResponse): unknown {
    const obj: any = {};
    message.replayProtection !== undefined &&
      (obj.replayProtection = message.replayProtection
        ? ReplayProtection.toJSON(message.replayProtection)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetReplayProtectionResponse>
  ): QueryGetReplayProtectionResponse {
    const message = {
      ...baseQueryGetReplayProtectionResponse,
    } as QueryGetReplayProtectionResponse;
    if (
      object.replayProtection !== undefined &&
      object.replayProtection !== null
    ) {
      message.replayProtection = ReplayProtection.fromPartial(
        object.replayProtection
      );
    } else {
      message.replayProtection = undefined;
    }
    return message;
  },
};

const baseQueryAllReplayProtectionRequest: object = {};

export const QueryAllReplayProtectionRequest = {
  encode(
    message: QueryAllReplayProtectionRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryAllReplayProtectionRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllReplayProtectionRequest,
    } as QueryAllReplayProtectionRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllReplayProtectionRequest {
    const message = {
      ...baseQueryAllReplayProtectionRequest,
    } as QueryAllReplayProtectionRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllReplayProtectionRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageRequest.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllReplayProtectionRequest>
  ): QueryAllReplayProtectionRequest {
    const message = {
      ...baseQueryAllReplayProtectionRequest,
    } as QueryAllReplayProtectionRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryAllReplayProtectionResponse: object = {};

export const QueryAllReplayProtectionResponse = {
  encode(
    message: QueryAllReplayProtectionResponse,
    writer: Writer = Writer.create()
  ): Writer {
    for (const v of message.replayProtection) {
      ReplayProtection.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(
        message.pagination,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryAllReplayProtectionResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllReplayProtectionResponse,
    } as QueryAllReplayProtectionResponse;
    message.replayProtection = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.replayProtection.push(
            ReplayProtection.decode(reader, reader.uint32())
          );
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllReplayProtectionResponse {
    const message = {
      ...baseQueryAllReplayProtectionResponse,
    } as QueryAllReplayProtectionResponse;
    message.replayProtection = [];
    if (
      object.replayProtection !== undefined &&
      object.replayProtection !== null
    ) {
      for (const e of object.replayProtection) {
        message.replayProtection.push(ReplayProtection.fromJSON(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllReplayProtectionResponse): unknown {
    const obj: any = {};
    if (message.replayProtection) {
      obj.replayProtection = message.replayProtection.map((e) =>
        e ? ReplayProtection.toJSON(e) : undefined
      );
    } else {
      obj.replayProtection = [];
    }
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageResponse.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllReplayProtectionResponse>
  ): QueryAllReplayProtectionResponse {
    const message = {
      ...baseQueryAllReplayProtectionResponse,
    } as QueryAllReplayProtectionResponse;
    message.replayProtection = [];
    if (
      object.replayProtection !== undefined &&
      object.replayProtection !== null
    ) {
      for (const e of object.replayProtection) {
        message.replayProtection.push(ReplayProtection.fromPartial(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryGetSequenceCounterRequest: object = { index: "" };

export const QueryGetSequenceCounterRequest = {
  encode(
    message: QueryGetSequenceCounterRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.index !== "") {
      writer.uint32(10).string(message.index);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryGetSequenceCounterRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryGetSequenceCounterRequest,
    } as QueryGetSequenceCounterRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.index = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetSequenceCounterRequest {
    const message = {
      ...baseQueryGetSequenceCounterRequest,
    } as QueryGetSequenceCounterRequest;
    if (object.index !== undefined && object.index !== null) {
      message.index = String(object.index);
    } else {
      message.index = "";
    }
    return message;
  },

  toJSON(message: QueryGetSequenceCounterRequest): unknown {
    const obj: any = {};
    message.index !== undefined && (obj.index = message.index);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetSequenceCounterRequest>
  ): QueryGetSequenceCounterRequest {
    const message = {
      ...baseQueryGetSequenceCounterRequest,
    } as QueryGetSequenceCounterRequest;
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    } else {
      message.index = "";
    }
    return message;
  },
};

const baseQueryGetSequenceCounterResponse: object = {};

export const QueryGetSequenceCounterResponse = {
  encode(
    message: QueryGetSequenceCounterResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.sequenceCounter !== undefined) {
      SequenceCounter.encode(
        message.sequenceCounter,
        writer.uint32(10).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryGetSequenceCounterResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryGetSequenceCounterResponse,
    } as QueryGetSequenceCounterResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sequenceCounter = SequenceCounter.decode(
            reader,
            reader.uint32()
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetSequenceCounterResponse {
    const message = {
      ...baseQueryGetSequenceCounterResponse,
    } as QueryGetSequenceCounterResponse;
    if (
      object.sequenceCounter !== undefined &&
      object.sequenceCounter !== null
    ) {
      message.sequenceCounter = SequenceCounter.fromJSON(
        object.sequenceCounter
      );
    } else {
      message.sequenceCounter = undefined;
    }
    return message;
  },

  toJSON(message: QueryGetSequenceCounterResponse): unknown {
    const obj: any = {};
    message.sequenceCounter !== undefined &&
      (obj.sequenceCounter = message.sequenceCounter
        ? SequenceCounter.toJSON(message.sequenceCounter)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetSequenceCounterResponse>
  ): QueryGetSequenceCounterResponse {
    const message = {
      ...baseQueryGetSequenceCounterResponse,
    } as QueryGetSequenceCounterResponse;
    if (
      object.sequenceCounter !== undefined &&
      object.sequenceCounter !== null
    ) {
      message.sequenceCounter = SequenceCounter.fromPartial(
        object.sequenceCounter
      );
    } else {
      message.sequenceCounter = undefined;
    }
    return message;
  },
};

const baseQueryAllSequenceCounterRequest: object = {};

export const QueryAllSequenceCounterRequest = {
  encode(
    message: QueryAllSequenceCounterRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryAllSequenceCounterRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllSequenceCounterRequest,
    } as QueryAllSequenceCounterRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllSequenceCounterRequest {
    const message = {
      ...baseQueryAllSequenceCounterRequest,
    } as QueryAllSequenceCounterRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllSequenceCounterRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageRequest.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllSequenceCounterRequest>
  ): QueryAllSequenceCounterRequest {
    const message = {
      ...baseQueryAllSequenceCounterRequest,
    } as QueryAllSequenceCounterRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryAllSequenceCounterResponse: object = {};

export const QueryAllSequenceCounterResponse = {
  encode(
    message: QueryAllSequenceCounterResponse,
    writer: Writer = Writer.create()
  ): Writer {
    for (const v of message.sequenceCounter) {
      SequenceCounter.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(
        message.pagination,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryAllSequenceCounterResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllSequenceCounterResponse,
    } as QueryAllSequenceCounterResponse;
    message.sequenceCounter = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sequenceCounter.push(
            SequenceCounter.decode(reader, reader.uint32())
          );
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllSequenceCounterResponse {
    const message = {
      ...baseQueryAllSequenceCounterResponse,
    } as QueryAllSequenceCounterResponse;
    message.sequenceCounter = [];
    if (
      object.sequenceCounter !== undefined &&
      object.sequenceCounter !== null
    ) {
      for (const e of object.sequenceCounter) {
        message.sequenceCounter.push(SequenceCounter.fromJSON(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllSequenceCounterResponse): unknown {
    const obj: any = {};
    if (message.sequenceCounter) {
      obj.sequenceCounter = message.sequenceCounter.map((e) =>
        e ? SequenceCounter.toJSON(e) : undefined
      );
    } else {
      obj.sequenceCounter = [];
    }
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageResponse.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllSequenceCounterResponse>
  ): QueryAllSequenceCounterResponse {
    const message = {
      ...baseQueryAllSequenceCounterResponse,
    } as QueryAllSequenceCounterResponse;
    message.sequenceCounter = [];
    if (
      object.sequenceCounter !== undefined &&
      object.sequenceCounter !== null
    ) {
      for (const e of object.sequenceCounter) {
        message.sequenceCounter.push(SequenceCounter.fromPartial(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryGetConsensusPhylaxSetIndexRequest: object = {};

export const QueryGetConsensusPhylaxSetIndexRequest = {
  encode(
    _: QueryGetConsensusPhylaxSetIndexRequest,
    writer: Writer = Writer.create()
  ): Writer {
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryGetConsensusPhylaxSetIndexRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryGetConsensusPhylaxSetIndexRequest,
    } as QueryGetConsensusPhylaxSetIndexRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): QueryGetConsensusPhylaxSetIndexRequest {
    const message = {
      ...baseQueryGetConsensusPhylaxSetIndexRequest,
    } as QueryGetConsensusPhylaxSetIndexRequest;
    return message;
  },

  toJSON(_: QueryGetConsensusPhylaxSetIndexRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<QueryGetConsensusPhylaxSetIndexRequest>
  ): QueryGetConsensusPhylaxSetIndexRequest {
    const message = {
      ...baseQueryGetConsensusPhylaxSetIndexRequest,
    } as QueryGetConsensusPhylaxSetIndexRequest;
    return message;
  },
};

const baseQueryGetConsensusPhylaxSetIndexResponse: object = {};

export const QueryGetConsensusPhylaxSetIndexResponse = {
  encode(
    message: QueryGetConsensusPhylaxSetIndexResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.ConsensusPhylaxSetIndex !== undefined) {
      ConsensusPhylaxSetIndex.encode(
        message.ConsensusPhylaxSetIndex,
        writer.uint32(10).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryGetConsensusPhylaxSetIndexResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryGetConsensusPhylaxSetIndexResponse,
    } as QueryGetConsensusPhylaxSetIndexResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.ConsensusPhylaxSetIndex = ConsensusPhylaxSetIndex.decode(
            reader,
            reader.uint32()
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetConsensusPhylaxSetIndexResponse {
    const message = {
      ...baseQueryGetConsensusPhylaxSetIndexResponse,
    } as QueryGetConsensusPhylaxSetIndexResponse;
    if (
      object.ConsensusPhylaxSetIndex !== undefined &&
      object.ConsensusPhylaxSetIndex !== null
    ) {
      message.ConsensusPhylaxSetIndex = ConsensusPhylaxSetIndex.fromJSON(
        object.ConsensusPhylaxSetIndex
      );
    } else {
      message.ConsensusPhylaxSetIndex = undefined;
    }
    return message;
  },

  toJSON(message: QueryGetConsensusPhylaxSetIndexResponse): unknown {
    const obj: any = {};
    message.ConsensusPhylaxSetIndex !== undefined &&
      (obj.ConsensusPhylaxSetIndex = message.ConsensusPhylaxSetIndex
        ? ConsensusPhylaxSetIndex.toJSON(message.ConsensusPhylaxSetIndex)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetConsensusPhylaxSetIndexResponse>
  ): QueryGetConsensusPhylaxSetIndexResponse {
    const message = {
      ...baseQueryGetConsensusPhylaxSetIndexResponse,
    } as QueryGetConsensusPhylaxSetIndexResponse;
    if (
      object.ConsensusPhylaxSetIndex !== undefined &&
      object.ConsensusPhylaxSetIndex !== null
    ) {
      message.ConsensusPhylaxSetIndex = ConsensusPhylaxSetIndex.fromPartial(
        object.ConsensusPhylaxSetIndex
      );
    } else {
      message.ConsensusPhylaxSetIndex = undefined;
    }
    return message;
  },
};

const baseQueryGetPhylaxValidatorRequest: object = {};

export const QueryGetPhylaxValidatorRequest = {
  encode(
    message: QueryGetPhylaxValidatorRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.guardianKey.length !== 0) {
      writer.uint32(10).bytes(message.guardianKey);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryGetPhylaxValidatorRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryGetPhylaxValidatorRequest,
    } as QueryGetPhylaxValidatorRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.guardianKey = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetPhylaxValidatorRequest {
    const message = {
      ...baseQueryGetPhylaxValidatorRequest,
    } as QueryGetPhylaxValidatorRequest;
    if (object.guardianKey !== undefined && object.guardianKey !== null) {
      message.guardianKey = bytesFromBase64(object.guardianKey);
    }
    return message;
  },

  toJSON(message: QueryGetPhylaxValidatorRequest): unknown {
    const obj: any = {};
    message.guardianKey !== undefined &&
      (obj.guardianKey = base64FromBytes(
        message.guardianKey !== undefined
          ? message.guardianKey
          : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetPhylaxValidatorRequest>
  ): QueryGetPhylaxValidatorRequest {
    const message = {
      ...baseQueryGetPhylaxValidatorRequest,
    } as QueryGetPhylaxValidatorRequest;
    if (object.guardianKey !== undefined && object.guardianKey !== null) {
      message.guardianKey = object.guardianKey;
    } else {
      message.guardianKey = new Uint8Array();
    }
    return message;
  },
};

const baseQueryGetPhylaxValidatorResponse: object = {};

export const QueryGetPhylaxValidatorResponse = {
  encode(
    message: QueryGetPhylaxValidatorResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.guardianValidator !== undefined) {
      PhylaxValidator.encode(
        message.guardianValidator,
        writer.uint32(10).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryGetPhylaxValidatorResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryGetPhylaxValidatorResponse,
    } as QueryGetPhylaxValidatorResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.guardianValidator = PhylaxValidator.decode(
            reader,
            reader.uint32()
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetPhylaxValidatorResponse {
    const message = {
      ...baseQueryGetPhylaxValidatorResponse,
    } as QueryGetPhylaxValidatorResponse;
    if (
      object.guardianValidator !== undefined &&
      object.guardianValidator !== null
    ) {
      message.guardianValidator = PhylaxValidator.fromJSON(
        object.guardianValidator
      );
    } else {
      message.guardianValidator = undefined;
    }
    return message;
  },

  toJSON(message: QueryGetPhylaxValidatorResponse): unknown {
    const obj: any = {};
    message.guardianValidator !== undefined &&
      (obj.guardianValidator = message.guardianValidator
        ? PhylaxValidator.toJSON(message.guardianValidator)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetPhylaxValidatorResponse>
  ): QueryGetPhylaxValidatorResponse {
    const message = {
      ...baseQueryGetPhylaxValidatorResponse,
    } as QueryGetPhylaxValidatorResponse;
    if (
      object.guardianValidator !== undefined &&
      object.guardianValidator !== null
    ) {
      message.guardianValidator = PhylaxValidator.fromPartial(
        object.guardianValidator
      );
    } else {
      message.guardianValidator = undefined;
    }
    return message;
  },
};

const baseQueryAllPhylaxValidatorRequest: object = {};

export const QueryAllPhylaxValidatorRequest = {
  encode(
    message: QueryAllPhylaxValidatorRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryAllPhylaxValidatorRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllPhylaxValidatorRequest,
    } as QueryAllPhylaxValidatorRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllPhylaxValidatorRequest {
    const message = {
      ...baseQueryAllPhylaxValidatorRequest,
    } as QueryAllPhylaxValidatorRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllPhylaxValidatorRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageRequest.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllPhylaxValidatorRequest>
  ): QueryAllPhylaxValidatorRequest {
    const message = {
      ...baseQueryAllPhylaxValidatorRequest,
    } as QueryAllPhylaxValidatorRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryAllPhylaxValidatorResponse: object = {};

export const QueryAllPhylaxValidatorResponse = {
  encode(
    message: QueryAllPhylaxValidatorResponse,
    writer: Writer = Writer.create()
  ): Writer {
    for (const v of message.guardianValidator) {
      PhylaxValidator.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(
        message.pagination,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryAllPhylaxValidatorResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllPhylaxValidatorResponse,
    } as QueryAllPhylaxValidatorResponse;
    message.guardianValidator = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.guardianValidator.push(
            PhylaxValidator.decode(reader, reader.uint32())
          );
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllPhylaxValidatorResponse {
    const message = {
      ...baseQueryAllPhylaxValidatorResponse,
    } as QueryAllPhylaxValidatorResponse;
    message.guardianValidator = [];
    if (
      object.guardianValidator !== undefined &&
      object.guardianValidator !== null
    ) {
      for (const e of object.guardianValidator) {
        message.guardianValidator.push(PhylaxValidator.fromJSON(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllPhylaxValidatorResponse): unknown {
    const obj: any = {};
    if (message.guardianValidator) {
      obj.guardianValidator = message.guardianValidator.map((e) =>
        e ? PhylaxValidator.toJSON(e) : undefined
      );
    } else {
      obj.guardianValidator = [];
    }
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageResponse.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllPhylaxValidatorResponse>
  ): QueryAllPhylaxValidatorResponse {
    const message = {
      ...baseQueryAllPhylaxValidatorResponse,
    } as QueryAllPhylaxValidatorResponse;
    message.guardianValidator = [];
    if (
      object.guardianValidator !== undefined &&
      object.guardianValidator !== null
    ) {
      for (const e of object.guardianValidator) {
        message.guardianValidator.push(PhylaxValidator.fromPartial(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryLatestPhylaxSetIndexRequest: object = {};

export const QueryLatestPhylaxSetIndexRequest = {
  encode(
    _: QueryLatestPhylaxSetIndexRequest,
    writer: Writer = Writer.create()
  ): Writer {
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryLatestPhylaxSetIndexRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryLatestPhylaxSetIndexRequest,
    } as QueryLatestPhylaxSetIndexRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): QueryLatestPhylaxSetIndexRequest {
    const message = {
      ...baseQueryLatestPhylaxSetIndexRequest,
    } as QueryLatestPhylaxSetIndexRequest;
    return message;
  },

  toJSON(_: QueryLatestPhylaxSetIndexRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<QueryLatestPhylaxSetIndexRequest>
  ): QueryLatestPhylaxSetIndexRequest {
    const message = {
      ...baseQueryLatestPhylaxSetIndexRequest,
    } as QueryLatestPhylaxSetIndexRequest;
    return message;
  },
};

const baseQueryLatestPhylaxSetIndexResponse: object = {
  latestPhylaxSetIndex: 0,
};

export const QueryLatestPhylaxSetIndexResponse = {
  encode(
    message: QueryLatestPhylaxSetIndexResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.latestPhylaxSetIndex !== 0) {
      writer.uint32(8).uint32(message.latestPhylaxSetIndex);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryLatestPhylaxSetIndexResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryLatestPhylaxSetIndexResponse,
    } as QueryLatestPhylaxSetIndexResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.latestPhylaxSetIndex = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryLatestPhylaxSetIndexResponse {
    const message = {
      ...baseQueryLatestPhylaxSetIndexResponse,
    } as QueryLatestPhylaxSetIndexResponse;
    if (
      object.latestPhylaxSetIndex !== undefined &&
      object.latestPhylaxSetIndex !== null
    ) {
      message.latestPhylaxSetIndex = Number(object.latestPhylaxSetIndex);
    } else {
      message.latestPhylaxSetIndex = 0;
    }
    return message;
  },

  toJSON(message: QueryLatestPhylaxSetIndexResponse): unknown {
    const obj: any = {};
    message.latestPhylaxSetIndex !== undefined &&
      (obj.latestPhylaxSetIndex = message.latestPhylaxSetIndex);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryLatestPhylaxSetIndexResponse>
  ): QueryLatestPhylaxSetIndexResponse {
    const message = {
      ...baseQueryLatestPhylaxSetIndexResponse,
    } as QueryLatestPhylaxSetIndexResponse;
    if (
      object.latestPhylaxSetIndex !== undefined &&
      object.latestPhylaxSetIndex !== null
    ) {
      message.latestPhylaxSetIndex = object.latestPhylaxSetIndex;
    } else {
      message.latestPhylaxSetIndex = 0;
    }
    return message;
  },
};

/** Query defines the gRPC querier service. */
export interface Query {
  /** Queries a guardianSet by index. */
  PhylaxSet(
    request: QueryGetPhylaxSetRequest
  ): Promise<QueryGetPhylaxSetResponse>;
  /** Queries a list of guardianSet items. */
  PhylaxSetAll(
    request: QueryAllPhylaxSetRequest
  ): Promise<QueryAllPhylaxSetResponse>;
  /** Queries a config by index. */
  Config(request: QueryGetConfigRequest): Promise<QueryGetConfigResponse>;
  /** Queries a replayProtection by index. */
  ReplayProtection(
    request: QueryGetReplayProtectionRequest
  ): Promise<QueryGetReplayProtectionResponse>;
  /** Queries a list of replayProtection items. */
  ReplayProtectionAll(
    request: QueryAllReplayProtectionRequest
  ): Promise<QueryAllReplayProtectionResponse>;
  /** Queries a sequenceCounter by index. */
  SequenceCounter(
    request: QueryGetSequenceCounterRequest
  ): Promise<QueryGetSequenceCounterResponse>;
  /** Queries a list of sequenceCounter items. */
  SequenceCounterAll(
    request: QueryAllSequenceCounterRequest
  ): Promise<QueryAllSequenceCounterResponse>;
  /** Queries a ConsensusPhylaxSetIndex by index. */
  ConsensusPhylaxSetIndex(
    request: QueryGetConsensusPhylaxSetIndexRequest
  ): Promise<QueryGetConsensusPhylaxSetIndexResponse>;
  /** Queries a PhylaxValidator by index. */
  PhylaxValidator(
    request: QueryGetPhylaxValidatorRequest
  ): Promise<QueryGetPhylaxValidatorResponse>;
  /** Queries a list of PhylaxValidator items. */
  PhylaxValidatorAll(
    request: QueryAllPhylaxValidatorRequest
  ): Promise<QueryAllPhylaxValidatorResponse>;
  /** Queries a list of LatestPhylaxSetIndex items. */
  LatestPhylaxSetIndex(
    request: QueryLatestPhylaxSetIndexRequest
  ): Promise<QueryLatestPhylaxSetIndexResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
  }
  PhylaxSet(
    request: QueryGetPhylaxSetRequest
  ): Promise<QueryGetPhylaxSetResponse> {
    const data = QueryGetPhylaxSetRequest.encode(request).finish();
    const promise = this.rpc.request(
      "wormhole_foundation.deltachain.wormhole.Query",
      "PhylaxSet",
      data
    );
    return promise.then((data) =>
      QueryGetPhylaxSetResponse.decode(new Reader(data))
    );
  }

  PhylaxSetAll(
    request: QueryAllPhylaxSetRequest
  ): Promise<QueryAllPhylaxSetResponse> {
    const data = QueryAllPhylaxSetRequest.encode(request).finish();
    const promise = this.rpc.request(
      "wormhole_foundation.deltachain.wormhole.Query",
      "PhylaxSetAll",
      data
    );
    return promise.then((data) =>
      QueryAllPhylaxSetResponse.decode(new Reader(data))
    );
  }

  Config(request: QueryGetConfigRequest): Promise<QueryGetConfigResponse> {
    const data = QueryGetConfigRequest.encode(request).finish();
    const promise = this.rpc.request(
      "wormhole_foundation.deltachain.wormhole.Query",
      "Config",
      data
    );
    return promise.then((data) =>
      QueryGetConfigResponse.decode(new Reader(data))
    );
  }

  ReplayProtection(
    request: QueryGetReplayProtectionRequest
  ): Promise<QueryGetReplayProtectionResponse> {
    const data = QueryGetReplayProtectionRequest.encode(request).finish();
    const promise = this.rpc.request(
      "wormhole_foundation.deltachain.wormhole.Query",
      "ReplayProtection",
      data
    );
    return promise.then((data) =>
      QueryGetReplayProtectionResponse.decode(new Reader(data))
    );
  }

  ReplayProtectionAll(
    request: QueryAllReplayProtectionRequest
  ): Promise<QueryAllReplayProtectionResponse> {
    const data = QueryAllReplayProtectionRequest.encode(request).finish();
    const promise = this.rpc.request(
      "wormhole_foundation.deltachain.wormhole.Query",
      "ReplayProtectionAll",
      data
    );
    return promise.then((data) =>
      QueryAllReplayProtectionResponse.decode(new Reader(data))
    );
  }

  SequenceCounter(
    request: QueryGetSequenceCounterRequest
  ): Promise<QueryGetSequenceCounterResponse> {
    const data = QueryGetSequenceCounterRequest.encode(request).finish();
    const promise = this.rpc.request(
      "wormhole_foundation.deltachain.wormhole.Query",
      "SequenceCounter",
      data
    );
    return promise.then((data) =>
      QueryGetSequenceCounterResponse.decode(new Reader(data))
    );
  }

  SequenceCounterAll(
    request: QueryAllSequenceCounterRequest
  ): Promise<QueryAllSequenceCounterResponse> {
    const data = QueryAllSequenceCounterRequest.encode(request).finish();
    const promise = this.rpc.request(
      "wormhole_foundation.deltachain.wormhole.Query",
      "SequenceCounterAll",
      data
    );
    return promise.then((data) =>
      QueryAllSequenceCounterResponse.decode(new Reader(data))
    );
  }

  ConsensusPhylaxSetIndex(
    request: QueryGetConsensusPhylaxSetIndexRequest
  ): Promise<QueryGetConsensusPhylaxSetIndexResponse> {
    const data = QueryGetConsensusPhylaxSetIndexRequest.encode(
      request
    ).finish();
    const promise = this.rpc.request(
      "wormhole_foundation.deltachain.wormhole.Query",
      "ConsensusPhylaxSetIndex",
      data
    );
    return promise.then((data) =>
      QueryGetConsensusPhylaxSetIndexResponse.decode(new Reader(data))
    );
  }

  PhylaxValidator(
    request: QueryGetPhylaxValidatorRequest
  ): Promise<QueryGetPhylaxValidatorResponse> {
    const data = QueryGetPhylaxValidatorRequest.encode(request).finish();
    const promise = this.rpc.request(
      "wormhole_foundation.deltachain.wormhole.Query",
      "PhylaxValidator",
      data
    );
    return promise.then((data) =>
      QueryGetPhylaxValidatorResponse.decode(new Reader(data))
    );
  }

  PhylaxValidatorAll(
    request: QueryAllPhylaxValidatorRequest
  ): Promise<QueryAllPhylaxValidatorResponse> {
    const data = QueryAllPhylaxValidatorRequest.encode(request).finish();
    const promise = this.rpc.request(
      "wormhole_foundation.deltachain.wormhole.Query",
      "PhylaxValidatorAll",
      data
    );
    return promise.then((data) =>
      QueryAllPhylaxValidatorResponse.decode(new Reader(data))
    );
  }

  LatestPhylaxSetIndex(
    request: QueryLatestPhylaxSetIndexRequest
  ): Promise<QueryLatestPhylaxSetIndexResponse> {
    const data = QueryLatestPhylaxSetIndexRequest.encode(request).finish();
    const promise = this.rpc.request(
      "wormhole_foundation.deltachain.wormhole.Query",
      "LatestPhylaxSetIndex",
      data
    );
    return promise.then((data) =>
      QueryLatestPhylaxSetIndexResponse.decode(new Reader(data))
    );
  }
}

interface Rpc {
  request(
    service: string,
    method: string,
    data: Uint8Array
  ): Promise<Uint8Array>;
}

declare var self: any | undefined;
declare var window: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") return globalThis;
  if (typeof self !== "undefined") return self;
  if (typeof window !== "undefined") return window;
  if (typeof global !== "undefined") return global;
  throw "Unable to locate global object";
})();

const atob: (b64: string) => string =
  globalThis.atob ||
  ((b64) => globalThis.Buffer.from(b64, "base64").toString("binary"));
function bytesFromBase64(b64: string): Uint8Array {
  const bin = atob(b64);
  const arr = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; ++i) {
    arr[i] = bin.charCodeAt(i);
  }
  return arr;
}

const btoa: (bin: string) => string =
  globalThis.btoa ||
  ((bin) => globalThis.Buffer.from(bin, "binary").toString("base64"));
function base64FromBytes(arr: Uint8Array): string {
  const bin: string[] = [];
  for (let i = 0; i < arr.byteLength; ++i) {
    bin.push(String.fromCharCode(arr[i]));
  }
  return btoa(bin.join(""));
}

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
