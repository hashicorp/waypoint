/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Token, NewTokenResponse } from 'waypoint-pb';
import { Response } from 'miragejs';
import { RouteHandler } from '../types';

function createToken(): Token {
  let token = new Token();
  let invite = new Token.Invite();
  token.setAccessorId(
    '3fwxJnSh32T9skH8NqseY8wuLQQynN6cnBYUCLTSxRJ6QCqLdEtUTY4hHjdDyHUiAarZC7WH1gZWypmQg8noi8ELfJxRe5131BFQWW3wzGW'
  );
  token.setInvite(invite);
  return token;
}

export function create(this: RouteHandler): Response {
  let resp = new NewTokenResponse();
  resp.setToken(createToken().getAccessorId_asB64());
  return this.serialize(resp, 'application');
}
