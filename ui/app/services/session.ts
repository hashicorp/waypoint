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
    window.localStorage.waypointAuthToken = value;
  }

  async removeToken() {
    window.localStorage.removeItem('waypointAuthToken');
  }
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    session: SessionService;
  }
}
