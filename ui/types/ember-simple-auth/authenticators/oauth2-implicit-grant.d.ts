/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

interface parseResponseObject {
  authuser: string;
  prompt: string;
  scope: string;
  code: string;
  state: string;
}
declare module 'ember-simple-auth/authenticators/oauth2-implicit-grant' {
  export default class OAuth2ImplicitGrantAuthenticator {}
  export function parseResponse(url: string): parseResponseObject;
}
