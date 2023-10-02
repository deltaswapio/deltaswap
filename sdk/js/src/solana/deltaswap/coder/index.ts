import { Coder, Idl } from "@project-serum/anchor";
import { DeltaswapAccountsCoder } from "./accounts";
import { DeltaswapEventsCoder } from "./events";
import { DeltaswapInstructionCoder } from "./instruction";
import { DeltaswapStateCoder } from "./state";
import { DeltaswapTypesCoder } from "./types";

export { DeltaswapInstruction } from "./instruction";

export class DeltaswapCoder implements Coder {
  readonly instruction: DeltaswapInstructionCoder;
  readonly accounts: DeltaswapAccountsCoder;
  readonly state: DeltaswapStateCoder;
  readonly events: DeltaswapEventsCoder;
  readonly types: DeltaswapTypesCoder;

  constructor(idl: Idl) {
    this.instruction = new DeltaswapInstructionCoder(idl);
    this.accounts = new DeltaswapAccountsCoder(idl);
    this.state = new DeltaswapStateCoder(idl);
    this.events = new DeltaswapEventsCoder(idl);
    this.types = new DeltaswapTypesCoder(idl);
  }
}
