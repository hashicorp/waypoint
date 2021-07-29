import DeploymentDetail from './deployment';
import { GetDeploymentRequest, Deployment, Ref } from 'waypoint-pb';

interface DeploymentIdModelParams {
  deployment_id: string;
}

export default class DeploymentIdDetail extends DeploymentDetail {
  renderTemplate() {
    this.render('workspace/projects/project/app/deployment', {
      into: 'workspace/projects/project',
    });
  }

  async model(params: DeploymentIdModelParams): Promise<Deployment.AsObject> {
    let ref = new Ref.Operation();
    ref.setId(params.deployment_id);
    let req = new GetDeploymentRequest();
    req.setRef(ref);

    let resp = await this.api.client.getDeployment(req, this.api.WithMeta());
    let deploy: Deployment = resp;
    return deploy.toObject();
  }
}
