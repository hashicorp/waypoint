/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// package: google.api
// file: google/api/context.proto

import * as jspb from "google-protobuf";

export class Context extends jspb.Message {
  clearRulesList(): void;
  getRulesList(): Array<ContextRule>;
  setRulesList(value: Array<ContextRule>): void;
  addRules(value?: ContextRule, index?: number): ContextRule;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Context.AsObject;
  static toObject(includeInstance: boolean, msg: Context): Context.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Context, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Context;
  static deserializeBinaryFromReader(message: Context, reader: jspb.BinaryReader): Context;
}

export namespace Context {
  export type AsObject = {
    rulesList: Array<ContextRule.AsObject>,
  }
}

export class ContextRule extends jspb.Message {
  getSelector(): string;
  setSelector(value: string): void;

  clearRequestedList(): void;
  getRequestedList(): Array<string>;
  setRequestedList(value: Array<string>): void;
  addRequested(value: string, index?: number): string;

  clearProvidedList(): void;
  getProvidedList(): Array<string>;
  setProvidedList(value: Array<string>): void;
  addProvided(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ContextRule.AsObject;
  static toObject(includeInstance: boolean, msg: ContextRule): ContextRule.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ContextRule, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ContextRule;
  static deserializeBinaryFromReader(message: ContextRule, reader: jspb.BinaryReader): ContextRule;
}

export namespace ContextRule {
  export type AsObject = {
    selector: string,
    requestedList: Array<string>,
    providedList: Array<string>,
  }
}

