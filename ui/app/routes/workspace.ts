import Route from '@ember/routing/route';
import { Ref } from 'waypoint-pb';

export type Params = { workspace_id: string };
export type Model = Ref.Workspace.AsObject;

export default class Workspace extends Route {
  async model(params: Params): Promise<Model> {
    // Workspace "id" which is a name, based on URL param
    let ws = new Ref.Workspace();
    ws.setWorkspace(params.workspace_id);

    return ws.toObject();
  }
}
