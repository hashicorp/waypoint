// package: google.api
// file: google/api/consumer.proto

import * as jspb from "google-protobuf";

export class ProjectProperties extends jspb.Message {
  clearPropertiesList(): void;
  getPropertiesList(): Array<Property>;
  setPropertiesList(value: Array<Property>): void;
  addProperties(value?: Property, index?: number): Property;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectProperties.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectProperties): ProjectProperties.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProjectProperties, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectProperties;
  static deserializeBinaryFromReader(message: ProjectProperties, reader: jspb.BinaryReader): ProjectProperties;
}

export namespace ProjectProperties {
  export type AsObject = {
    propertiesList: Array<Property.AsObject>,
  }
}

export class Property extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getType(): Property.PropertyTypeMap[keyof Property.PropertyTypeMap];
  setType(value: Property.PropertyTypeMap[keyof Property.PropertyTypeMap]): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Property.AsObject;
  static toObject(includeInstance: boolean, msg: Property): Property.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Property, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Property;
  static deserializeBinaryFromReader(message: Property, reader: jspb.BinaryReader): Property;
}

export namespace Property {
  export type AsObject = {
    name: string,
    type: Property.PropertyTypeMap[keyof Property.PropertyTypeMap],
    description: string,
  }

  export interface PropertyTypeMap {
    UNSPECIFIED: 0;
    INT64: 1;
    BOOL: 2;
    STRING: 3;
    DOUBLE: 4;
  }

  export const PropertyType: PropertyTypeMap;
}

