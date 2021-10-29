import { reject, resolve } from 'rsvp';

import BaseAuthenticator from 'ember-simple-auth/authenticators/base';
import classic from 'ember-classic-decorator';

interface SessionData {
  token?: string;
}

@classic
export default class TokenAuthenticator extends BaseAuthenticator {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
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
