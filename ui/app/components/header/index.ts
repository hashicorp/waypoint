import Component from '@glimmer/component';
import { SessionService } from 'ember-simple-auth/services/session';
import { inject as service } from '@ember/service';
export default class Header extends Component {
  @service session!: SessionService;

  get canLogout(): boolean {
    return this.session.isAuthenticated;
  }
}
