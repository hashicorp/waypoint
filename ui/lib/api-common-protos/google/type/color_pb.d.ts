// package: google.type
// file: google/type/color.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class Color extends jspb.Message {
  getRed(): number;
  setRed(value: number): void;

  getGreen(): number;
  setGreen(value: number): void;

  getBlue(): number;
  setBlue(value: number): void;

  hasAlpha(): boolean;
  clearAlpha(): void;
  getAlpha(): google_protobuf_wrappers_pb.FloatValue | undefined;
  setAlpha(value?: google_protobuf_wrappers_pb.FloatValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Color.AsObject;
  static toObject(includeInstance: boolean, msg: Color): Color.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Color, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Color;
  static deserializeBinaryFromReader(message: Color, reader: jspb.BinaryReader): Color;
}

export namespace Color {
  export type AsObject = {
    red: number,
    green: number,
    blue: number,
    alpha?: google_protobuf_wrappers_pb.FloatValue.AsObject,
  }
}

