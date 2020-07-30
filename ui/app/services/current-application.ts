import Service, { inject as service } from '@ember/service';
import { Application, Ref } from 'waypoint-pb';
import ApiService from 'waypoint/services/api';

export default class CurrentApplicationService extends Service {
  @service api!: ApiService;

  application?: Application;
  ref?: Ref.Application;

  setRef(ref: Ref.Application) {
    this.ref = ref;
  }

  setApplication(application: Application) {
    this.application = application;
  }
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    currentApplication: CurrentApplicationService;
  }
}
