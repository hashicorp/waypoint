/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// package: google.api
// file: google/api/quota.proto

import * as jspb from "google-protobuf";
import * as google_api_annotations_pb from "../../google/api/annotations_pb";

export class Quota extends jspb.Message {
  clearLimitsList(): void;
  getLimitsList(): Array<QuotaLimit>;
  setLimitsList(value: Array<QuotaLimit>): void;
  addLimits(value?: QuotaLimit, index?: number): QuotaLimit;

  clearMetricRulesList(): void;
  getMetricRulesList(): Array<MetricRule>;
  setMetricRulesList(value: Array<MetricRule>): void;
  addMetricRules(value?: MetricRule, index?: number): MetricRule;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Quota.AsObject;
  static toObject(includeInstance: boolean, msg: Quota): Quota.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Quota, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Quota;
  static deserializeBinaryFromReader(message: Quota, reader: jspb.BinaryReader): Quota;
}

export namespace Quota {
  export type AsObject = {
    limitsList: Array<QuotaLimit.AsObject>,
    metricRulesList: Array<MetricRule.AsObject>,
  }
}

export class MetricRule extends jspb.Message {
  getSelector(): string;
  setSelector(value: string): void;

  getMetricCostsMap(): jspb.Map<string, number>;
  clearMetricCostsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MetricRule.AsObject;
  static toObject(includeInstance: boolean, msg: MetricRule): MetricRule.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MetricRule, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MetricRule;
  static deserializeBinaryFromReader(message: MetricRule, reader: jspb.BinaryReader): MetricRule;
}

export namespace MetricRule {
  export type AsObject = {
    selector: string,
    metricCostsMap: Array<[string, number]>,
  }
}

export class QuotaLimit extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getDefaultLimit(): number;
  setDefaultLimit(value: number): void;

  getMaxLimit(): number;
  setMaxLimit(value: number): void;

  getFreeTier(): number;
  setFreeTier(value: number): void;

  getDuration(): string;
  setDuration(value: string): void;

  getMetric(): string;
  setMetric(value: string): void;

  getUnit(): string;
  setUnit(value: string): void;

  getValuesMap(): jspb.Map<string, number>;
  clearValuesMap(): void;
  getDisplayName(): string;
  setDisplayName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QuotaLimit.AsObject;
  static toObject(includeInstance: boolean, msg: QuotaLimit): QuotaLimit.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: QuotaLimit, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QuotaLimit;
  static deserializeBinaryFromReader(message: QuotaLimit, reader: jspb.BinaryReader): QuotaLimit;
}

export namespace QuotaLimit {
  export type AsObject = {
    name: string,
    description: string,
    defaultLimit: number,
    maxLimit: number,
    freeTier: number,
    duration: string,
    metric: string,
    unit: string,
    valuesMap: Array<[string, number]>,
    displayName: string,
  }
}

