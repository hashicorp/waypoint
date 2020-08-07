import FakeXMLHttpRequest from 'fake-xml-http-request';
import faker from './faker';

import { Component } from 'waypoint-pb';

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
export function logRequestConsole(verb: string, path: string, request: FakeXMLHttpRequest) {
  let url = request.responseURL.split('/');
  if (url.length >= 5) {
    console.groupCollapsed(`Mock: ${url[2]} ${url[3]}/${url[4]}`);
  }
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

  console.groupCollapsed('Response (base64 decoded)');
  console.log(loggedResponse);
  console.groupEnd();

  console.groupCollapsed('Request (base64 decoded)');
  console.log(loggedRequest);
  console.groupEnd();

  console.groupCollapsed('Request (raw)');
  console.log(request);
  console.groupEnd();

  console.groupEnd();
}
