import Base from 'ember-simple-auth/authenticators/base';
import classic from 'ember-classic-decorator';

@classic
export default class TokenAuthenticator extends Base {
  restore(data) {

  }

  authenticate(token) {
    debugger;
  }

  invalidate(data) {

  }
};
