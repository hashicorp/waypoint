import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { GetDeploymentRequest, Deployment } from 'waypoint-pb';
import CurrentApplicationService from 'waypoint/services/current-application';
import CurrentWorkspaceService from 'waypoint/services/current-workspace';

interface DeploymentModelParams {
  deployment_id: string;
}

export default class DeploymentDetail extends Route {
  @service api!: ApiService;
  @service currentApplication!: CurrentApplicationService;
  @service currentWorkspace!: CurrentWorkspaceService;

  async model(params: DeploymentModelParams) {
    var req = new GetDeploymentRequest();
    req.setDeploymentId(params.deployment_id);

    var resp = await this.api.client.getDeployment(req, {});
    let deploy: Deployment = resp;
    return deploy.toObject();
  }
}
