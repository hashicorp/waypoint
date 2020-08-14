import Service, { inject as service } from '@ember/service';
import ApiService from './api';

export default class SessionService extends Service {
  @service api!: ApiService;

  get authConfigured(): Boolean {
    return Boolean(this.token);
  }

  get token(): string {
    return window.localStorage.waypointAuthToken;
  }

  async setToken(value: string) {
    // todo(pearkes): validate this either locally (format) or remotely (rpc)
    if (value == null) {
      window.localStorage.removeItem('waypointAuthToken');
    } else {
      window.localStorage.waypointAuthToken = value;
    }
  }
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    session: SessionService;
  }
}
