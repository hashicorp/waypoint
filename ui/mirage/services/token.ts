/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Token, NewTokenResponse } from 'waypoint-pb';
import { Response } from 'miragejs';
import { RouteHandler } from '../types';

function createToken(): Token {
  let token = new Token();
  token.setAccessorId(
    'xjLoe9b2j2jTYLpM5vAv1Z6JrUW438HabDQ7fvyBzWCLozD6L2oBWRE8G2zk62V1UzcrcGf1LvmwbuQYyAFwRq63n3996WHsqrbyb8XLXfAqDbCePNX96Fkt'
  );
  return token;
}

export function create(this: RouteHandler): Response {
  let resp = new NewTokenResponse();
  resp.setToken(createToken().getAccessorId_asB64());
  return this.serialize(resp, 'application');
}
