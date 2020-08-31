import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Ref } from 'waypoint-pb';

interface WSModelParams {
  workspace_id: string;
}

export default class Workspace extends Route {
  async model(params: WSModelParams): Promise<Ref.Workspace.AsObject> {
    // Workspace "id" which is a name, based on URL param
    let ws = new Ref.Workspace();
    ws.setWorkspace(params.workspace_id);

    return ws.toObject();
  }
}
