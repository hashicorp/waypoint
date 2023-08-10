/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { reject, resolve } from 'rsvp';

import BaseAuthenticator from 'ember-simple-auth/authenticators/base';
import classic from 'ember-classic-decorator';

interface SessionData {
  token?: string;
}

@classic
export default class TokenAuthenticator extends BaseAuthenticator {
  restore(data: SessionData): Promise<SessionData> {
    if (data.token) {
      return resolve(data);
    } else {
      return reject();
    }
  }

  authenticate(token: string): Promise<SessionData> {
    if (token !== '') {
      return resolve({ token: token });
    } else {
      return reject();
    }
  }
}
