import Service, { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import ApiService from './api';

export default class SessionService extends Service {
  @service api!: ApiService;
  @tracked authConfigured = false;

  get token(): string {
    return window.localStorage.waypointAuthToken;
  }

  async setToken(value: string) {
    window.localStorage.waypointAuthToken = value;
    this.authConfigured = true;
  }

  async removeToken() {
    window.localStorage.removeItem('waypointAuthToken');
    this.authConfigured = false;
  }
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    session: SessionService;
  }
}
