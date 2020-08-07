import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { ListDeploymentsRequest, ListDeploymentsResponse } from 'waypoint-pb';
import CurrentApplicationService from 'waypoint/services/current-application';
import CurrentWorkspaceService from 'waypoint/services/current-workspace';

export default class Deployments extends Route {
  @service api!: ApiService;
  @service currentApplication!: CurrentApplicationService;
  @service currentWorkspace!: CurrentWorkspaceService;

  async model() {
    var req = new ListDeploymentsRequest();
    req.setApplication(this.currentApplication.ref);
    req.setWorkspace(this.currentWorkspace.ref);

    var resp = await this.api.client.listDeployments(req, {});
    let deployResp: ListDeploymentsResponse = resp;

    return deployResp.getDeploymentsList().map((b) => b.toObject());
  }
}
