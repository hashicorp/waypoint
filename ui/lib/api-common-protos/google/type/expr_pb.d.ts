// package: google.type
// file: google/type/expr.proto

import * as jspb from "google-protobuf";

export class Expr extends jspb.Message {
  getExpression(): string;
  setExpression(value: string): void;

  getTitle(): string;
  setTitle(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getLocation(): string;
  setLocation(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Expr.AsObject;
  static toObject(includeInstance: boolean, msg: Expr): Expr.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Expr, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Expr;
  static deserializeBinaryFromReader(message: Expr, reader: jspb.BinaryReader): Expr;
}

export namespace Expr {
  export type AsObject = {
    expression: string,
    title: string,
    description: string,
    location: string,
  }
}

