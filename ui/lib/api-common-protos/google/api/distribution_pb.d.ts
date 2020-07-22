// package: google.api
// file: google/api/distribution.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

export class Distribution extends jspb.Message {
  getCount(): number;
  setCount(value: number): void;

  getMean(): number;
  setMean(value: number): void;

  getSumOfSquaredDeviation(): number;
  setSumOfSquaredDeviation(value: number): void;

  hasRange(): boolean;
  clearRange(): void;
  getRange(): Distribution.Range | undefined;
  setRange(value?: Distribution.Range): void;

  hasBucketOptions(): boolean;
  clearBucketOptions(): void;
  getBucketOptions(): Distribution.BucketOptions | undefined;
  setBucketOptions(value?: Distribution.BucketOptions): void;

  clearBucketCountsList(): void;
  getBucketCountsList(): Array<number>;
  setBucketCountsList(value: Array<number>): void;
  addBucketCounts(value: number, index?: number): number;

  clearExemplarsList(): void;
  getExemplarsList(): Array<Distribution.Exemplar>;
  setExemplarsList(value: Array<Distribution.Exemplar>): void;
  addExemplars(value?: Distribution.Exemplar, index?: number): Distribution.Exemplar;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Distribution.AsObject;
  static toObject(includeInstance: boolean, msg: Distribution): Distribution.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Distribution, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Distribution;
  static deserializeBinaryFromReader(message: Distribution, reader: jspb.BinaryReader): Distribution;
}

export namespace Distribution {
  export type AsObject = {
    count: number,
    mean: number,
    sumOfSquaredDeviation: number,
    range?: Distribution.Range.AsObject,
    bucketOptions?: Distribution.BucketOptions.AsObject,
    bucketCountsList: Array<number>,
    exemplarsList: Array<Distribution.Exemplar.AsObject>,
  }

  export class Range extends jspb.Message {
    getMin(): number;
    setMin(value: number): void;

    getMax(): number;
    setMax(value: number): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Range.AsObject;
    static toObject(includeInstance: boolean, msg: Range): Range.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Range, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Range;
    static deserializeBinaryFromReader(message: Range, reader: jspb.BinaryReader): Range;
  }

  export namespace Range {
    export type AsObject = {
      min: number,
      max: number,
    }
  }

  export class BucketOptions extends jspb.Message {
    hasLinearBuckets(): boolean;
    clearLinearBuckets(): void;
    getLinearBuckets(): Distribution.BucketOptions.Linear | undefined;
    setLinearBuckets(value?: Distribution.BucketOptions.Linear): void;

    hasExponentialBuckets(): boolean;
    clearExponentialBuckets(): void;
    getExponentialBuckets(): Distribution.BucketOptions.Exponential | undefined;
    setExponentialBuckets(value?: Distribution.BucketOptions.Exponential): void;

    hasExplicitBuckets(): boolean;
    clearExplicitBuckets(): void;
    getExplicitBuckets(): Distribution.BucketOptions.Explicit | undefined;
    setExplicitBuckets(value?: Distribution.BucketOptions.Explicit): void;

    getOptionsCase(): BucketOptions.OptionsCase;
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): BucketOptions.AsObject;
    static toObject(includeInstance: boolean, msg: BucketOptions): BucketOptions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: BucketOptions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): BucketOptions;
    static deserializeBinaryFromReader(message: BucketOptions, reader: jspb.BinaryReader): BucketOptions;
  }

  export namespace BucketOptions {
    export type AsObject = {
      linearBuckets?: Distribution.BucketOptions.Linear.AsObject,
      exponentialBuckets?: Distribution.BucketOptions.Exponential.AsObject,
      explicitBuckets?: Distribution.BucketOptions.Explicit.AsObject,
    }

    export class Linear extends jspb.Message {
      getNumFiniteBuckets(): number;
      setNumFiniteBuckets(value: number): void;

      getWidth(): number;
      setWidth(value: number): void;

      getOffset(): number;
      setOffset(value: number): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Linear.AsObject;
      static toObject(includeInstance: boolean, msg: Linear): Linear.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: Linear, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Linear;
      static deserializeBinaryFromReader(message: Linear, reader: jspb.BinaryReader): Linear;
    }

    export namespace Linear {
      export type AsObject = {
        numFiniteBuckets: number,
        width: number,
        offset: number,
      }
    }

    export class Exponential extends jspb.Message {
      getNumFiniteBuckets(): number;
      setNumFiniteBuckets(value: number): void;

      getGrowthFactor(): number;
      setGrowthFactor(value: number): void;

      getScale(): number;
      setScale(value: number): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Exponential.AsObject;
      static toObject(includeInstance: boolean, msg: Exponential): Exponential.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: Exponential, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Exponential;
      static deserializeBinaryFromReader(message: Exponential, reader: jspb.BinaryReader): Exponential;
    }

    export namespace Exponential {
      export type AsObject = {
        numFiniteBuckets: number,
        growthFactor: number,
        scale: number,
      }
    }

    export class Explicit extends jspb.Message {
      clearBoundsList(): void;
      getBoundsList(): Array<number>;
      setBoundsList(value: Array<number>): void;
      addBounds(value: number, index?: number): number;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Explicit.AsObject;
      static toObject(includeInstance: boolean, msg: Explicit): Explicit.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: Explicit, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Explicit;
      static deserializeBinaryFromReader(message: Explicit, reader: jspb.BinaryReader): Explicit;
    }

    export namespace Explicit {
      export type AsObject = {
        boundsList: Array<number>,
      }
    }

    export enum OptionsCase {
      OPTIONS_NOT_SET = 0,
      LINEAR_BUCKETS = 1,
      EXPONENTIAL_BUCKETS = 2,
      EXPLICIT_BUCKETS = 3,
    }
  }

  export class Exemplar extends jspb.Message {
    getValue(): number;
    setValue(value: number): void;

    hasTimestamp(): boolean;
    clearTimestamp(): void;
    getTimestamp(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setTimestamp(value?: google_protobuf_timestamp_pb.Timestamp): void;

    clearAttachmentsList(): void;
    getAttachmentsList(): Array<google_protobuf_any_pb.Any>;
    setAttachmentsList(value: Array<google_protobuf_any_pb.Any>): void;
    addAttachments(value?: google_protobuf_any_pb.Any, index?: number): google_protobuf_any_pb.Any;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Exemplar.AsObject;
    static toObject(includeInstance: boolean, msg: Exemplar): Exemplar.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Exemplar, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Exemplar;
    static deserializeBinaryFromReader(message: Exemplar, reader: jspb.BinaryReader): Exemplar;
  }

  export namespace Exemplar {
    export type AsObject = {
      value: number,
      timestamp?: google_protobuf_timestamp_pb.Timestamp.AsObject,
      attachmentsList: Array<google_protobuf_any_pb.Any.AsObject>,
    }
  }
}

