import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { ListBuildsRequest, Ref, ListBuildsResponse } from 'waypoint-pb';

export default class Builds extends Route {
  @service api!: ApiService;

  async model() {
    let ws: Ref.Workspace = await this.modelFor('workspace').ref;
    let project: Ref.Project = await this.modelFor('workspace.project').ref;
    let app: Ref.Application = await this.modelFor('workspace.project.app').ref;

    var req = new ListBuildsRequest();
    req.setApplication(app);
    req.setWorkspace(ws);

    var resp = await this.api.client.listBuilds(req, {});
    let buildResp: ListBuildsResponse = resp;

    return buildResp.getBuildsList().map((b) => b.toObject());
  }
}
