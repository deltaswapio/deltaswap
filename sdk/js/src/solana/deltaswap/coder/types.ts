import { Idl, TypesCoder } from "@project-serum/anchor";

export class DeltaswapTypesCoder implements TypesCoder {
  constructor(_idl: Idl) {}

  encode<T = any>(_name: string, _type: T): Buffer {
    throw new Error("Deltaswap program does not have user-defined types");
  }
  decode<T = any>(_name: string, _typeData: Buffer): T {
    throw new Error("Deltaswap program does not have user-defined types");
  }
}
