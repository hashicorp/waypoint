/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import FakeXMLHttpRequest from 'fake-xml-http-request';
import faker from './faker';
import { Component } from 'waypoint-pb';
import { Status } from 'waypoint-pb';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { subMinutes } from 'date-fns';

/**
 * Creates a random status object in a random state and random amount of time
 * passed
 */
export function sequenceRandom(): number {
  return faker.random.number({
    min: 5,
    max: 800,
  });
}

/**
 * Creates a random status object in a random state and random amount of time
 * passed
 */
export function statusRandom(): Status {
  let status = new Status();
  status.setState(Status.State.SUCCESS);

  let startTs = new Timestamp();
  let endTs = new Timestamp();

  let minutes = faker.random.number({
    min: 5,
    max: 40,
  });

  startTs.setSeconds(Math.floor(subMinutes(new Date(), minutes).getTime() / 1000));
  endTs.setSeconds(Math.floor(new Date().getTime() / 1000));

  status.setStartTime(startTs);
  status.setCompleteTime(endTs);

  return status;
}

/**
 * Generates a random fake ID roughly matching the Waypoint
 * internal format.
 */
export function fakeId(): string {
  return faker.random.alphaNumeric(24).toUpperCase();
}

/**
 * Known component types and names to Waypoint
 * for mocking and development.
 */
export const componentOptions = {
  1: ['pack', 'docker'],
  2: ['docker'],
  3: ['google-cloud-run', 'kubernetes', 'docker'],
  4: ['google-cloud-run', 'kubernetes', 'docker'],
};

/**
 * Returns an object that has a random component type and
 * component name.
 */
export function fakeComponentForKind(kind: Component.Type): string {
  return componentOptions[kind][Math.floor(Math.random() * componentOptions[kind].length)];
}

/**
 * Logs a FakeXMLHttpRequest to the console, trying to decode base64 encoded
 * request and response bodies for debugging. This is only partially helpful and
 * dependent on the readability of the data being sent.
 */
export function logRequestConsole(verb: string, path: string, request: FakeXMLHttpRequest): void {
  console.groupCollapsed(`Mock: ${verb} ${path}`);

  let { requestBody, responseText } = request;
  let loggedRequest: string, loggedResponse: string;

  try {
    loggedRequest = atob(requestBody);
  } catch (e) {
    loggedRequest = requestBody;
  }

  try {
    loggedResponse = atob(responseText);
  } catch (e) {
    loggedResponse = responseText;
  }

  console.groupCollapsed('Request (raw)');
  console.log(request);
  console.groupEnd();

  console.groupCollapsed('Request (protobuf wire format)');
  console.log(loggedRequest);
  console.groupEnd();

  console.groupCollapsed('Response (protobuf wire format)');
  console.log(loggedResponse);
  console.groupEnd();

  console.groupEnd();
}

export function dateToTimestamp(date: Date): Timestamp {
  let result = new Timestamp();

  result.setSeconds(Math.floor(date.valueOf() / 1000));

  return result;
}
