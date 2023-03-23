/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// package: google.api
// file: google/api/backend.proto

import * as jspb from "google-protobuf";

export class Backend extends jspb.Message {
  clearRulesList(): void;
  getRulesList(): Array<BackendRule>;
  setRulesList(value: Array<BackendRule>): void;
  addRules(value?: BackendRule, index?: number): BackendRule;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Backend.AsObject;
  static toObject(includeInstance: boolean, msg: Backend): Backend.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Backend, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Backend;
  static deserializeBinaryFromReader(message: Backend, reader: jspb.BinaryReader): Backend;
}

export namespace Backend {
  export type AsObject = {
    rulesList: Array<BackendRule.AsObject>,
  }
}

export class BackendRule extends jspb.Message {
  getSelector(): string;
  setSelector(value: string): void;

  getAddress(): string;
  setAddress(value: string): void;

  getDeadline(): number;
  setDeadline(value: number): void;

  getMinDeadline(): number;
  setMinDeadline(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BackendRule.AsObject;
  static toObject(includeInstance: boolean, msg: BackendRule): BackendRule.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BackendRule, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BackendRule;
  static deserializeBinaryFromReader(message: BackendRule, reader: jspb.BinaryReader): BackendRule;
}

export namespace BackendRule {
  export type AsObject = {
    selector: string,
    address: string,
    deadline: number,
    minDeadline: number,
  }
}

