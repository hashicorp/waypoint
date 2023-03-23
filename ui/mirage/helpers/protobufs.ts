/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import * as protobuf from 'google-protobuf';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

/**
 * Encodes the given protobuf message to base64-encoded wire format.
 *
 * @param msg the protobuf message you want to encode
 * @returns msg in base64-encoded protobuf wire format
 */
export function encode(msg: protobuf.Message): string {
  let resp = msg || new Empty();
  let serialized = resp.serializeBinary();
  let len = serialized.length;
  let bytesArray = [0, 0, 0, 0];
  let payload = new Uint8Array(5 + len);

  for (let i = 3; i >= 0; i--) {
    bytesArray[i] = len % 256;
    len = len >>> 8;
  }

  payload.set(new Uint8Array(bytesArray), 1);
  payload.set(serialized, 5);

  let result = btoa(String.fromCharCode(...payload));

  return result;
}

/**
 * Decodes base64-encoded wire format to the given message type.
 *
 * @param type the type of protobuf message you’re expecting to decode
 * @param data the base64-encoded protobuf wire format
 * @returns reified protobuf message
 */
export function decode<T extends protobuf.Message>(
  type: { deserializeBinary(msg: Uint8Array): T },
  data: string
): T {
  let wireBase64 = data;
  let wireAscii = atob(wireBase64);
  let wireChars = [...wireAscii];
  let msgChars = wireChars.slice(5);
  let msgCharCodes = msgChars.map((s) => s.charCodeAt(0));
  let msgBinary = Uint8Array.from(msgCharCodes);
  let result = type.deserializeBinary(msgBinary);

  return result;
}
