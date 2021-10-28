import Service, { inject as service } from '@ember/service';

import ApiService from './api';
import { tracked } from '@glimmer/tracking';

export default class OldSessionService extends Service {
  @service api!: ApiService;
  @tracked authConfigured: boolean;

  constructor(...args: ConstructorParameters<typeof Service>) {
    super(...args);

    this.authConfigured = false;
    if (this.token) {
      this.authConfigured = true;
    }
  }

  get token(): string {
    return window.localStorage.waypointAuthToken;
  }

  async setToken(value: string): Promise<void> {
    this.authConfigured = true;
    window.localStorage.waypointAuthToken = value;
  }

  async removeToken(): Promise<void> {
    this.authConfigured = false;
    window.localStorage.removeItem('waypointAuthToken');
  }
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    session: SessionService;
  }
}
