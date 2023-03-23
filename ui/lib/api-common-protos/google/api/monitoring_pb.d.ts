/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// package: google.api
// file: google/api/monitoring.proto

import * as jspb from "google-protobuf";
import * as google_api_annotations_pb from "../../google/api/annotations_pb";

export class Monitoring extends jspb.Message {
  clearProducerDestinationsList(): void;
  getProducerDestinationsList(): Array<Monitoring.MonitoringDestination>;
  setProducerDestinationsList(value: Array<Monitoring.MonitoringDestination>): void;
  addProducerDestinations(value?: Monitoring.MonitoringDestination, index?: number): Monitoring.MonitoringDestination;

  clearConsumerDestinationsList(): void;
  getConsumerDestinationsList(): Array<Monitoring.MonitoringDestination>;
  setConsumerDestinationsList(value: Array<Monitoring.MonitoringDestination>): void;
  addConsumerDestinations(value?: Monitoring.MonitoringDestination, index?: number): Monitoring.MonitoringDestination;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Monitoring.AsObject;
  static toObject(includeInstance: boolean, msg: Monitoring): Monitoring.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Monitoring, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Monitoring;
  static deserializeBinaryFromReader(message: Monitoring, reader: jspb.BinaryReader): Monitoring;
}

export namespace Monitoring {
  export type AsObject = {
    producerDestinationsList: Array<Monitoring.MonitoringDestination.AsObject>,
    consumerDestinationsList: Array<Monitoring.MonitoringDestination.AsObject>,
  }

  export class MonitoringDestination extends jspb.Message {
    getMonitoredResource(): string;
    setMonitoredResource(value: string): void;

    clearMetricsList(): void;
    getMetricsList(): Array<string>;
    setMetricsList(value: Array<string>): void;
    addMetrics(value: string, index?: number): string;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): MonitoringDestination.AsObject;
    static toObject(includeInstance: boolean, msg: MonitoringDestination): MonitoringDestination.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: MonitoringDestination, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): MonitoringDestination;
    static deserializeBinaryFromReader(message: MonitoringDestination, reader: jspb.BinaryReader): MonitoringDestination;
  }

  export namespace MonitoringDestination {
    export type AsObject = {
      monitoredResource: string,
      metricsList: Array<string>,
    }
  }
}

