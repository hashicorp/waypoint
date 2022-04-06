import SessionService from 'ember-simple-auth/services/session';
import classic from 'ember-classic-decorator';

@classic
class WaypointSessionService extends SessionService {
  handleInvalidation() {
    this._super(...arguments);
    this.set('data.workspace', undefined);
  }
}

export default WaypointSessionService;
