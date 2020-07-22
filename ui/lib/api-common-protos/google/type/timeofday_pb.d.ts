// package: google.type
// file: google/type/timeofday.proto

import * as jspb from "google-protobuf";

export class TimeOfDay extends jspb.Message {
  getHours(): number;
  setHours(value: number): void;

  getMinutes(): number;
  setMinutes(value: number): void;

  getSeconds(): number;
  setSeconds(value: number): void;

  getNanos(): number;
  setNanos(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TimeOfDay.AsObject;
  static toObject(includeInstance: boolean, msg: TimeOfDay): TimeOfDay.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TimeOfDay, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TimeOfDay;
  static deserializeBinaryFromReader(message: TimeOfDay, reader: jspb.BinaryReader): TimeOfDay;
}

export namespace TimeOfDay {
  export type AsObject = {
    hours: number,
    minutes: number,
    seconds: number,
    nanos: number,
  }
}

