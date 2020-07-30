import Service, { inject as service } from '@ember/service';
import { Ref } from 'waypoint-pb';
import ApiService from 'waypoint/services/api';

export default class CurrentWorkspaceService extends Service {
  @service api!: ApiService;

  ref?: Ref.Workspace;

  setRef(ref: Ref.Workspace) {
    this.ref = ref;
  }
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    currentWorkspace: CurrentWorkspaceService;
  }
}
