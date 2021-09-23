import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import { action } from '@ember/object';
import SessionService from 'waypoint/services/session';
import RouterService from '@ember/routing/router-service';

export default class LoginForm extends Component {
  @service session!: SessionService;
  @service router!: RouterService;

  token = '';

  @action
  async login(event?: Event): Promise<void> {
    event?.preventDefault();

    await this.session.setToken(this.token);
    this.router.transitionTo('workspaces');
  }
}
