import Service, { inject as service } from '@ember/service';
import { Ref } from 'waypoint-pb';
import ApiService from 'waypoint/services/api';
import { tracked } from '@glimmer/tracking';

export default class CurrentWorkspaceService extends Service {
  @service api!: ApiService;

  @tracked ref?: Ref.Workspace;
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    currentWorkspace: CurrentWorkspaceService;
  }
}
