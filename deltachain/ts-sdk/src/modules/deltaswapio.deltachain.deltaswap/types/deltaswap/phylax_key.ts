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

export interface PhylaxKey {
  key: Uint8Array;
}

const basePhylaxKey: object = {};

export const PhylaxKey = {
  encode(message: PhylaxKey, writer: Writer = Writer.create()): Writer {
    if (message.key.length !== 0) {
      writer.uint32(10).bytes(message.key);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): PhylaxKey {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...basePhylaxKey } as PhylaxKey;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.key = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PhylaxKey {
    const message = { ...basePhylaxKey } as PhylaxKey;
    if (object.key !== undefined && object.key !== null) {
      message.key = bytesFromBase64(object.key);
    }
    return message;
  },

  toJSON(message: PhylaxKey): unknown {
    const obj: any = {};
    message.key !== undefined &&
      (obj.key = base64FromBytes(
        message.key !== undefined ? message.key : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(object: DeepPartial<PhylaxKey>): PhylaxKey {
    const message = { ...basePhylaxKey } as PhylaxKey;
    if (object.key !== undefined && object.key !== null) {
      message.key = object.key;
    } else {
      message.key = new Uint8Array();
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
