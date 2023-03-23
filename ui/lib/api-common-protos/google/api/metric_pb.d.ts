/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// package: google.api
// file: google/api/metric.proto

import * as jspb from "google-protobuf";
import * as google_api_label_pb from "../../google/api/label_pb";

export class MetricDescriptor extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getType(): string;
  setType(value: string): void;

  clearLabelsList(): void;
  getLabelsList(): Array<google_api_label_pb.LabelDescriptor>;
  setLabelsList(value: Array<google_api_label_pb.LabelDescriptor>): void;
  addLabels(value?: google_api_label_pb.LabelDescriptor, index?: number): google_api_label_pb.LabelDescriptor;

  getMetricKind(): MetricDescriptor.MetricKindMap[keyof MetricDescriptor.MetricKindMap];
  setMetricKind(value: MetricDescriptor.MetricKindMap[keyof MetricDescriptor.MetricKindMap]): void;

  getValueType(): MetricDescriptor.ValueTypeMap[keyof MetricDescriptor.ValueTypeMap];
  setValueType(value: MetricDescriptor.ValueTypeMap[keyof MetricDescriptor.ValueTypeMap]): void;

  getUnit(): string;
  setUnit(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getDisplayName(): string;
  setDisplayName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MetricDescriptor.AsObject;
  static toObject(includeInstance: boolean, msg: MetricDescriptor): MetricDescriptor.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MetricDescriptor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MetricDescriptor;
  static deserializeBinaryFromReader(message: MetricDescriptor, reader: jspb.BinaryReader): MetricDescriptor;
}

export namespace MetricDescriptor {
  export type AsObject = {
    name: string,
    type: string,
    labelsList: Array<google_api_label_pb.LabelDescriptor.AsObject>,
    metricKind: MetricDescriptor.MetricKindMap[keyof MetricDescriptor.MetricKindMap],
    valueType: MetricDescriptor.ValueTypeMap[keyof MetricDescriptor.ValueTypeMap],
    unit: string,
    description: string,
    displayName: string,
  }

  export interface MetricKindMap {
    METRIC_KIND_UNSPECIFIED: 0;
    GAUGE: 1;
    DELTA: 2;
    CUMULATIVE: 3;
  }

  export const MetricKind: MetricKindMap;

  export interface ValueTypeMap {
    VALUE_TYPE_UNSPECIFIED: 0;
    BOOL: 1;
    INT64: 2;
    DOUBLE: 3;
    STRING: 4;
    DISTRIBUTION: 5;
    MONEY: 6;
  }

  export const ValueType: ValueTypeMap;
}

export class Metric extends jspb.Message {
  getType(): string;
  setType(value: string): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Metric.AsObject;
  static toObject(includeInstance: boolean, msg: Metric): Metric.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Metric, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Metric;
  static deserializeBinaryFromReader(message: Metric, reader: jspb.BinaryReader): Metric;
}

export namespace Metric {
  export type AsObject = {
    type: string,
    labelsMap: Array<[string, string]>,
  }
}

