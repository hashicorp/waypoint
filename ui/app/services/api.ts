import Service from '@ember/service';
import { WaypointClient } from 'waypoint-client';
import SessionService from 'waypoint/services/session';
import { inject as service } from '@ember/service';
import { assign } from '@ember/polyfills';

export default class ApiService extends Service {
  @service session!: SessionService;
  meta = { authorization: this.session.token };

  client = new WaypointClient('https://localhost:1235', null, null);

  // Merges metadata with required metadata for the request
  WithMeta(meta?: any) {
    // In the future we may want additional metadata per-request so this
    // helper merges that per-request metadata supplied at the client request
    // with our authentication metadata
    return assign(this.meta, meta!).valueOf();
  }
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    api: ApiService;
  }
}
