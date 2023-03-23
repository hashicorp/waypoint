/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Response } from 'miragejs';
import { RouteHandler } from '../types';

import { ListOIDCAuthMethodsResponse } from 'waypoint-pb';

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function list(this: RouteHandler, schema: any): Response {
  let authMethods = schema.authMethods.all();
  let authMethodsProtos = authMethods.models?.map((model) => model?.toProtobuf());

  let resp = new ListOIDCAuthMethodsResponse();
  resp.setAuthMethodsList(authMethodsProtos);
  return this.serialize(resp, 'application');
}
