// package: google.api
// file: google/api/label.proto

import * as jspb from "google-protobuf";

export class LabelDescriptor extends jspb.Message {
  getKey(): string;
  setKey(value: string): void;

  getValueType(): LabelDescriptor.ValueTypeMap[keyof LabelDescriptor.ValueTypeMap];
  setValueType(value: LabelDescriptor.ValueTypeMap[keyof LabelDescriptor.ValueTypeMap]): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LabelDescriptor.AsObject;
  static toObject(includeInstance: boolean, msg: LabelDescriptor): LabelDescriptor.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LabelDescriptor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LabelDescriptor;
  static deserializeBinaryFromReader(message: LabelDescriptor, reader: jspb.BinaryReader): LabelDescriptor;
}

export namespace LabelDescriptor {
  export type AsObject = {
    key: string,
    valueType: LabelDescriptor.ValueTypeMap[keyof LabelDescriptor.ValueTypeMap],
    description: string,
  }

  export interface ValueTypeMap {
    STRING: 0;
    BOOL: 1;
    INT64: 2;
  }

  export const ValueType: ValueTypeMap;
}

