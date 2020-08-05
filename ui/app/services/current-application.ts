import Service, { inject as service } from '@ember/service';
import { Application, Ref } from 'waypoint-pb';
import ApiService from 'waypoint/services/api';
import { tracked } from '@glimmer/tracking';

export default class CurrentApplicationService extends Service {
  @service api!: ApiService;

  @tracked application?: Application;
  @tracked ref?: Ref.Application;
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    currentApplication: CurrentApplicationService;
  }
}
