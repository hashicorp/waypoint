import { reject, resolve } from 'rsvp';

import BaseAuthenticator from 'ember-simple-auth/authenticators/base';
import classic from 'ember-classic-decorator';

@classic
export default class TokenAuthenticator extends BaseAuthenticator {
  restore(data) {
    if (data.token) {
      return resolve(data);
    } else {
      return reject();
    }
  }

  authenticate(token) {
    if (token !== '') {
      return resolve({ token: token });
    } else {
      return reject();
    }
  }

  invalidate(data) {

  }
};
