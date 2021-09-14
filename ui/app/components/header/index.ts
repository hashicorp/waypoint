import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import SessionService from 'waypoint/services/session';

export default class Header extends Component {
  @service session!: SessionService;

  get canLogout(): boolean {
    return this.session.authConfigured;
  }
}
