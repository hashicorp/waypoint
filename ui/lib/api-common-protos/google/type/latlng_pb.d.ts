// package: google.type
// file: google/type/latlng.proto

import * as jspb from "google-protobuf";

export class LatLng extends jspb.Message {
  getLatitude(): number;
  setLatitude(value: number): void;

  getLongitude(): number;
  setLongitude(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LatLng.AsObject;
  static toObject(includeInstance: boolean, msg: LatLng): LatLng.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LatLng, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LatLng;
  static deserializeBinaryFromReader(message: LatLng, reader: jspb.BinaryReader): LatLng;
}

export namespace LatLng {
  export type AsObject = {
    latitude: number,
    longitude: number,
  }
}

