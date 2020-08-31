import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { GetDeploymentRequest, Deployment, Ref } from 'waypoint-pb';

interface DeploymentModelParams {
  deployment_id: string;
}

export default class DeploymentDetail extends Route {
  @service api!: ApiService;

  async model(params: DeploymentModelParams) {
    var ref = new Ref.Operation();
    ref.setId(params.deployment_id);
    var req = new GetDeploymentRequest();
    req.setRef(ref);

    var resp = await this.api.client.getDeployment(req, this.api.WithMeta());
    let deploy: Deployment = resp;
    return deploy.toObject();
  }
}
