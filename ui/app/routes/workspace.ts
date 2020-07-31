import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Ref } from 'waypoint-pb';
import CurrentWorkspaceService from 'waypoint/services/current-workspace';

interface WSModelParams {
  id: string;
}

export default class Workspace extends Route {
  @service api!: ApiService;
  @service currentWorkspace!: CurrentWorkspaceService;

  async model(params: WSModelParams) {
    // Workspace "id" which is a name, based on URL param
    let ws = new Ref.Workspace();
    ws.setWorkspace(params.id);

    // Set on service, note we do not have a Workspace
    this.currentWorkspace.setRef(ws);
    return ws.toObject();
  }
}
