/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// package: google.api
// file: google/api/resource.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_descriptor_pb from "google-protobuf/google/protobuf/descriptor_pb";

export class ResourceDescriptor extends jspb.Message {
  getType(): string;
  setType(value: string): void;

  clearPatternList(): void;
  getPatternList(): Array<string>;
  setPatternList(value: Array<string>): void;
  addPattern(value: string, index?: number): string;

  getNameField(): string;
  setNameField(value: string): void;

  getHistory(): ResourceDescriptor.HistoryMap[keyof ResourceDescriptor.HistoryMap];
  setHistory(value: ResourceDescriptor.HistoryMap[keyof ResourceDescriptor.HistoryMap]): void;

  getPlural(): string;
  setPlural(value: string): void;

  getSingular(): string;
  setSingular(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResourceDescriptor.AsObject;
  static toObject(includeInstance: boolean, msg: ResourceDescriptor): ResourceDescriptor.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResourceDescriptor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResourceDescriptor;
  static deserializeBinaryFromReader(message: ResourceDescriptor, reader: jspb.BinaryReader): ResourceDescriptor;
}

export namespace ResourceDescriptor {
  export type AsObject = {
    type: string,
    patternList: Array<string>,
    nameField: string,
    history: ResourceDescriptor.HistoryMap[keyof ResourceDescriptor.HistoryMap],
    plural: string,
    singular: string,
  }

  export interface HistoryMap {
    HISTORY_UNSPECIFIED: 0;
    ORIGINALLY_SINGLE_PATTERN: 1;
    FUTURE_MULTI_PATTERN: 2;
  }

  export const History: HistoryMap;
}

export class ResourceReference extends jspb.Message {
  getType(): string;
  setType(value: string): void;

  getChildType(): string;
  setChildType(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResourceReference.AsObject;
  static toObject(includeInstance: boolean, msg: ResourceReference): ResourceReference.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResourceReference, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResourceReference;
  static deserializeBinaryFromReader(message: ResourceReference, reader: jspb.BinaryReader): ResourceReference;
}

export namespace ResourceReference {
  export type AsObject = {
    type: string,
    childType: string,
  }
}

  export const resourceReference: jspb.ExtensionFieldInfo<ResourceReference>;

  export const resourceDefinition: jspb.ExtensionFieldInfo<ResourceDescriptor>;

  export const resource: jspb.ExtensionFieldInfo<ResourceDescriptor>;

