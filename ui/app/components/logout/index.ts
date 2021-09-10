import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import { action } from '@ember/object';
import SessionService from 'waypoint/services/session';
import RouterService from '@ember/routing/router-service';

export default class Logout extends Component {
  @service session!: SessionService;
  @service router!: RouterService;

  @action
  async logout(): Promise<void> {
    await this.session.removeToken();
    this.router.transitionTo('auth');
  }
}
