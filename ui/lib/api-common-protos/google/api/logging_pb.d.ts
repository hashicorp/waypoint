/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// package: google.api
// file: google/api/logging.proto

import * as jspb from "google-protobuf";
import * as google_api_annotations_pb from "../../google/api/annotations_pb";

export class Logging extends jspb.Message {
  clearProducerDestinationsList(): void;
  getProducerDestinationsList(): Array<Logging.LoggingDestination>;
  setProducerDestinationsList(value: Array<Logging.LoggingDestination>): void;
  addProducerDestinations(value?: Logging.LoggingDestination, index?: number): Logging.LoggingDestination;

  clearConsumerDestinationsList(): void;
  getConsumerDestinationsList(): Array<Logging.LoggingDestination>;
  setConsumerDestinationsList(value: Array<Logging.LoggingDestination>): void;
  addConsumerDestinations(value?: Logging.LoggingDestination, index?: number): Logging.LoggingDestination;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Logging.AsObject;
  static toObject(includeInstance: boolean, msg: Logging): Logging.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Logging, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Logging;
  static deserializeBinaryFromReader(message: Logging, reader: jspb.BinaryReader): Logging;
}

export namespace Logging {
  export type AsObject = {
    producerDestinationsList: Array<Logging.LoggingDestination.AsObject>,
    consumerDestinationsList: Array<Logging.LoggingDestination.AsObject>,
  }

  export class LoggingDestination extends jspb.Message {
    getMonitoredResource(): string;
    setMonitoredResource(value: string): void;

    clearLogsList(): void;
    getLogsList(): Array<string>;
    setLogsList(value: Array<string>): void;
    addLogs(value: string, index?: number): string;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): LoggingDestination.AsObject;
    static toObject(includeInstance: boolean, msg: LoggingDestination): LoggingDestination.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: LoggingDestination, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): LoggingDestination;
    static deserializeBinaryFromReader(message: LoggingDestination, reader: jspb.BinaryReader): LoggingDestination;
  }

  export namespace LoggingDestination {
    export type AsObject = {
      monitoredResource: string,
      logsList: Array<string>,
    }
  }
}

