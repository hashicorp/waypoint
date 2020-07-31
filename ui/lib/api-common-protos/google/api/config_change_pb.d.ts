// package: google.api
// file: google/api/config_change.proto

import * as jspb from "google-protobuf";

export class ConfigChange extends jspb.Message {
  getElement(): string;
  setElement(value: string): void;

  getOldValue(): string;
  setOldValue(value: string): void;

  getNewValue(): string;
  setNewValue(value: string): void;

  getChangeType(): ChangeTypeMap[keyof ChangeTypeMap];
  setChangeType(value: ChangeTypeMap[keyof ChangeTypeMap]): void;

  clearAdvicesList(): void;
  getAdvicesList(): Array<Advice>;
  setAdvicesList(value: Array<Advice>): void;
  addAdvices(value?: Advice, index?: number): Advice;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConfigChange.AsObject;
  static toObject(includeInstance: boolean, msg: ConfigChange): ConfigChange.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConfigChange, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConfigChange;
  static deserializeBinaryFromReader(message: ConfigChange, reader: jspb.BinaryReader): ConfigChange;
}

export namespace ConfigChange {
  export type AsObject = {
    element: string,
    oldValue: string,
    newValue: string,
    changeType: ChangeTypeMap[keyof ChangeTypeMap],
    advicesList: Array<Advice.AsObject>,
  }
}

export class Advice extends jspb.Message {
  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Advice.AsObject;
  static toObject(includeInstance: boolean, msg: Advice): Advice.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Advice, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Advice;
  static deserializeBinaryFromReader(message: Advice, reader: jspb.BinaryReader): Advice;
}

export namespace Advice {
  export type AsObject = {
    description: string,
  }
}

export interface ChangeTypeMap {
  CHANGE_TYPE_UNSPECIFIED: 0;
  ADDED: 1;
  REMOVED: 2;
  MODIFIED: 3;
}

export const ChangeType: ChangeTypeMap;

