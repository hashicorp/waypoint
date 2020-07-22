// package: google.type
// file: google/type/fraction.proto

import * as jspb from "google-protobuf";

export class Fraction extends jspb.Message {
  getNumerator(): number;
  setNumerator(value: number): void;

  getDenominator(): number;
  setDenominator(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Fraction.AsObject;
  static toObject(includeInstance: boolean, msg: Fraction): Fraction.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Fraction, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Fraction;
  static deserializeBinaryFromReader(message: Fraction, reader: jspb.BinaryReader): Fraction;
}

export namespace Fraction {
  export type AsObject = {
    numerator: number,
    denominator: number,
  }
}

