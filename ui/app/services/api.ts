import Service from '@ember/service';
import { WaypointClient } from 'waypoint-client';
import { assign } from '@ember/polyfills';

export default class ApiService extends Service {
  token =
    'bM152PWkXxfoy4vA51JFhR7LsKQez6x23oi2RDqYk8DPjnRdGjWtS6J3CTywJSaBPQX7wZAgV61bFMMLWoqvpjUfr1pL2sq9AcDGL';
  meta = { authorization: this.token };
  // opts = { unaryInterceptors: [new ExampleUnaryInterceptor()] };
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
