import { EventCoder, Event, Idl } from "@project-serum/anchor";
import { IdlEvent } from "../../anchor";

export class DeltaswapEventsCoder implements EventCoder {
  constructor(_idl: Idl) {}

  decode<E extends IdlEvent = IdlEvent, T = Record<string, string>>(
    _log: string
  ): Event<E, T> | null {
    throw new Error("Deltaswap program does not have events");
  }
}
