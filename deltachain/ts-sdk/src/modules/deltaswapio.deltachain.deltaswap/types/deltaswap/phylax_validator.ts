//@ts-nocheck
//@ts-nocheck
//@ts-nocheck
//@ts-nocheck
//@ts-nocheck
//@ts-nocheck
//@ts-nocheck
/* eslint-disable */
import { Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "deltaswapio.deltachain.deltaswap";

export interface PhylaxValidator {
  phylaxKey: Uint8Array;
  validatorAddr: Uint8Array;
}

const basePhylaxValidator: object = {};

export const PhylaxValidator = {
  encode(message: PhylaxValidator, writer: Writer = Writer.create()): Writer {
    if (message.phylaxKey.length !== 0) {
      writer.uint32(10).bytes(message.phylaxKey);
    }
    if (message.validatorAddr.length !== 0) {
      writer.uint32(18).bytes(message.validatorAddr);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): PhylaxValidator {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...basePhylaxValidator } as PhylaxValidator;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.phylaxKey = reader.bytes();
          break;
        case 2:
          message.validatorAddr = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PhylaxValidator {
    const message = { ...basePhylaxValidator } as PhylaxValidator;
    if (object.phylaxKey !== undefined && object.phylaxKey !== null) {
      message.phylaxKey = bytesFromBase64(object.phylaxKey);
    }
    if (object.validatorAddr !== undefined && object.validatorAddr !== null) {
      message.validatorAddr = bytesFromBase64(object.validatorAddr);
    }
    return message;
  },

  toJSON(message: PhylaxValidator): unknown {
    const obj: any = {};
    message.phylaxKey !== undefined &&
      (obj.phylaxKey = base64FromBytes(
        message.phylaxKey !== undefined
          ? message.phylaxKey
          : new Uint8Array()
      ));
    message.validatorAddr !== undefined &&
      (obj.validatorAddr = base64FromBytes(
        message.validatorAddr !== undefined
          ? message.validatorAddr
          : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(object: DeepPartial<PhylaxValidator>): PhylaxValidator {
    const message = { ...basePhylaxValidator } as PhylaxValidator;
    if (object.phylaxKey !== undefined && object.phylaxKey !== null) {
      message.phylaxKey = object.phylaxKey;
    } else {
      message.phylaxKey = new Uint8Array();
    }
    if (object.validatorAddr !== undefined && object.validatorAddr !== null) {
      message.validatorAddr = object.validatorAddr;
    } else {
      message.validatorAddr = new Uint8Array();
    }
    return message;
  },
};

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
