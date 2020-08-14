import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import SessionService from 'waypoint/services/session';
import Transition from '@ember/routing/-private/transition';

export default class Application extends Route {
  @service session!: SessionService;

  async beforeModel(transition: Transition) {
    await super.beforeModel(transition);

    console.log(this.session.authConfigured);

    if (!this.session.authConfigured) {
      this.transitionTo('auth');
    }
  }
}
