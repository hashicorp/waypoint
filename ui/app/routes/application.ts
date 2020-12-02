import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import SessionService from 'waypoint/services/session';
import Transition from '@ember/routing/-private/transition';
import { action } from '@ember/object';
import ApiService from 'waypoint/services/api';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';


const ErrInvalidToken = 'invalid authentication token';

export default class Application extends Route {
  @service session!: SessionService;
  @service api!: ApiService;

  async beforeModel(transition: Transition) {
    await super.beforeModel(transition);
    if (!this.session.authConfigured && !transition.to.name.startsWith('auth')) {
      this.transitionTo('auth');
    }
  }

  async model() {
    let resp = await this.api.client.getVersionInfo(new Empty(), this.api.WithMeta());
    let versionInfo = resp?.getInfo();
    return versionInfo?.toObject();
  }

  @action
  error(error: Error) {
    console.log(error);

    if (error.message.includes(ErrInvalidToken)) {
      this.session.removeToken();
      this.transitionTo('auth');
    }
    return true;
  }
}
