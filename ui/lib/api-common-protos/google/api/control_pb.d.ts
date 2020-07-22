// package: google.api
// file: google/api/control.proto

import * as jspb from "google-protobuf";

export class Control extends jspb.Message {
  getEnvironment(): string;
  setEnvironment(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Control.AsObject;
  static toObject(includeInstance: boolean, msg: Control): Control.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Control, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Control;
  static deserializeBinaryFromReader(message: Control, reader: jspb.BinaryReader): Control;
}

export namespace Control {
  export type AsObject = {
    environment: string,
  }
}

