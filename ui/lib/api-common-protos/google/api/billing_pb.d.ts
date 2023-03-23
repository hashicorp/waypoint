/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// package: google.api
// file: google/api/billing.proto

import * as jspb from "google-protobuf";
import * as google_api_annotations_pb from "../../google/api/annotations_pb";

export class Billing extends jspb.Message {
  clearConsumerDestinationsList(): void;
  getConsumerDestinationsList(): Array<Billing.BillingDestination>;
  setConsumerDestinationsList(value: Array<Billing.BillingDestination>): void;
  addConsumerDestinations(value?: Billing.BillingDestination, index?: number): Billing.BillingDestination;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Billing.AsObject;
  static toObject(includeInstance: boolean, msg: Billing): Billing.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Billing, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Billing;
  static deserializeBinaryFromReader(message: Billing, reader: jspb.BinaryReader): Billing;
}

export namespace Billing {
  export type AsObject = {
    consumerDestinationsList: Array<Billing.BillingDestination.AsObject>,
  }

  export class BillingDestination extends jspb.Message {
    getMonitoredResource(): string;
    setMonitoredResource(value: string): void;

    clearMetricsList(): void;
    getMetricsList(): Array<string>;
    setMetricsList(value: Array<string>): void;
    addMetrics(value: string, index?: number): string;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): BillingDestination.AsObject;
    static toObject(includeInstance: boolean, msg: BillingDestination): BillingDestination.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: BillingDestination, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): BillingDestination;
    static deserializeBinaryFromReader(message: BillingDestination, reader: jspb.BinaryReader): BillingDestination;
  }

  export namespace BillingDestination {
    export type AsObject = {
      monitoredResource: string,
      metricsList: Array<string>,
    }
  }
}

