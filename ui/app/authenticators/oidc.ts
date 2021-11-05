import OAuth2ImplicitGrantAuthenticator, {
  parseResponse as ESAparseResponse,
} from 'ember-simple-auth/authenticators/oauth2-implicit-grant';
import { reject, resolve } from 'rsvp';

import classic from 'ember-classic-decorator';

interface SessionData {
  token: string;
}

@classic
export default class OIDCAuthenticator extends OAuth2ImplicitGrantAuthenticator {
  restore(data: SessionData): Promise<SessionData> {
    if (data.token) {
      window.localStorage.removeItem('waypointOIDCAuthMethod');
      window.localStorage.removeItem('waypointOIDCNonce');
      return resolve(data);
    } else {
      return reject();
    }
  }

  authenticate(hash: SessionData): Promise<SessionData> {
    if (hash.token !== '') {
      return resolve(hash);
    } else {
      return reject();
    }
  }
}

export function parseResponse(args: string): unknown {
  return ESAparseResponse(args);
}
