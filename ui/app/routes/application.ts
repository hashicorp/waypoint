import Route from '@ember/routing/route';
import { SessionService } from 'ember-simple-auth/services/session';
import Transition from '@ember/routing/-private/transition';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';

const ErrInvalidToken = 'invalid authentication token';

export default class Application extends Route {
  @service session!: SessionService;

  async beforeModel(transition: Transition): Promise<void> {
    await super.beforeModel(transition);
    if (!this.session.isAuthenticated && !transition.to.name.startsWith('auth')) {
      this.transitionTo('auth');
    }
  }

  @action
  error(error: Error): boolean | void {
    console.log(error);

    if (error.message.includes(ErrInvalidToken)) {
      this.session.invalidate();
      this.transitionTo('auth');
    }
    return true;
  }
}
