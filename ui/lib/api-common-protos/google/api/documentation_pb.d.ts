/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// package: google.api
// file: google/api/documentation.proto

import * as jspb from "google-protobuf";

export class Documentation extends jspb.Message {
  getSummary(): string;
  setSummary(value: string): void;

  clearPagesList(): void;
  getPagesList(): Array<Page>;
  setPagesList(value: Array<Page>): void;
  addPages(value?: Page, index?: number): Page;

  clearRulesList(): void;
  getRulesList(): Array<DocumentationRule>;
  setRulesList(value: Array<DocumentationRule>): void;
  addRules(value?: DocumentationRule, index?: number): DocumentationRule;

  getDocumentationRootUrl(): string;
  setDocumentationRootUrl(value: string): void;

  getOverview(): string;
  setOverview(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Documentation.AsObject;
  static toObject(includeInstance: boolean, msg: Documentation): Documentation.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Documentation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Documentation;
  static deserializeBinaryFromReader(message: Documentation, reader: jspb.BinaryReader): Documentation;
}

export namespace Documentation {
  export type AsObject = {
    summary: string,
    pagesList: Array<Page.AsObject>,
    rulesList: Array<DocumentationRule.AsObject>,
    documentationRootUrl: string,
    overview: string,
  }
}

export class DocumentationRule extends jspb.Message {
  getSelector(): string;
  setSelector(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getDeprecationDescription(): string;
  setDeprecationDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DocumentationRule.AsObject;
  static toObject(includeInstance: boolean, msg: DocumentationRule): DocumentationRule.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DocumentationRule, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DocumentationRule;
  static deserializeBinaryFromReader(message: DocumentationRule, reader: jspb.BinaryReader): DocumentationRule;
}

export namespace DocumentationRule {
  export type AsObject = {
    selector: string,
    description: string,
    deprecationDescription: string,
  }
}

export class Page extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getContent(): string;
  setContent(value: string): void;

  clearSubpagesList(): void;
  getSubpagesList(): Array<Page>;
  setSubpagesList(value: Array<Page>): void;
  addSubpages(value?: Page, index?: number): Page;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Page.AsObject;
  static toObject(includeInstance: boolean, msg: Page): Page.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Page, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Page;
  static deserializeBinaryFromReader(message: Page, reader: jspb.BinaryReader): Page;
}

export namespace Page {
  export type AsObject = {
    name: string,
    content: string,
    subpagesList: Array<Page.AsObject>,
  }
}

