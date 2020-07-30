import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { ListBuildsRequest, ListBuildsResponse } from 'waypoint-pb';
import CurrentApplicationService from 'waypoint/services/current-application';
import CurrentWorkspaceService from 'waypoint/services/current-workspace';

export default class Builds extends Route {
  @service api!: ApiService;
  @service currentApplication!: CurrentApplicationService;
  @service currentWorkspace!: CurrentWorkspaceService;

  async model() {
    var req = new ListBuildsRequest();
    req.setApplication(this.currentApplication.ref);
    req.setWorkspace(this.currentWorkspace.ref);

    var resp = await this.api.client.listBuilds(req, {});
    let buildResp: ListBuildsResponse = resp;

    return buildResp.getBuildsList().map((b) => b.toObject());
  }
}
