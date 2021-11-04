import { reject, resolve } from 'rsvp';

import OAuth2ImplicitGrantAuthenticator from 'ember-simple-auth/authenticators/oauth2-implicit-grant';
import classic from 'ember-classic-decorator';

interface SessionData {
  token?: string;
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

  authenticate(token: string): Promise<SessionData> {
    if (token !== '') {
      return resolve({ token: token });
    } else {
      return reject();
    }
  }
}
