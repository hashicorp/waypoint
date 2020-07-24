import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Ref } from 'waypoint-pb';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

interface WSModelParams {
  id: string;
}

export default class Workspace extends Route {
  @service api!: ApiService;

  async model(params: WSModelParams) {
    let ws = new Ref.Workspace();

    // For now, assume a single default workspace
    ws.setWorkspace(params.id);
    let resp = await this.api.client.listProjects(new Empty(), {});
    let projects = resp.getProjectsList().map((p) => p.toObject());

    return {
      ref: ws as Ref.Workspace,
      workspace: ws.toObject(),
      projects: projects,
    };
  }
}
