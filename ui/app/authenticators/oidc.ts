/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import OAuth2ImplicitGrantAuthenticator, {
  parseResponse as ESAparseResponse,
} from 'ember-simple-auth/authenticators/oauth2-implicit-grant';
import { reject, resolve } from 'rsvp';

import classic from 'ember-classic-decorator';

interface SessionData {
  token: string;
}
interface parseResponseObject {
  authuser: string;
  prompt: string;
  scope: string;
  code: string;
  state: string;
}

@classic
export default class OIDCAuthenticator extends OAuth2ImplicitGrantAuthenticator {
  restore(data: SessionData): Promise<SessionData> {
    if (data.token) {
      return resolve(data);
    } else {
      return reject();
    }
  }

  authenticate(hash: SessionData): Promise<SessionData> {
    if (hash.token !== '') {
      this._cleanUpLocalStorage();
      return resolve(hash);
    } else {
      return reject();
    }
  }

  // Used to clean up OIDC information stored in LocalStorage
  // during the authentication flow
  _cleanUpLocalStorage(): void {
    window.localStorage.removeItem('waypointOIDCAuthMethod');
    window.localStorage.removeItem('waypointOIDCNonce');
  }
}

export function parseResponse(args: string): parseResponseObject {
  return ESAparseResponse(args);
}
