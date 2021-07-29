import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { GetDeploymentRequest, Deployment, Ref } from 'waypoint-pb';
import ApiService from 'waypoint/services/api';

type Params = { deployment_id: string };
type Model = Deployment.AsObject;

export default class DeploymentIdDetail extends Route {
  @service api!: ApiService;

  async model(params: Params): Promise<Model> {
    let req = new GetDeploymentRequest();
    let ref = new Ref.Operation();

    ref.setId(params.deployment_id);
    req.setRef(ref);

    let deployment = await this.api.client.getDeployment(req, this.api.WithMeta());

    return deployment.toObject();
  }

  redirect(model: Model): void {
    this.transitionTo('workspace.projects.project.app.deployment', model.sequence);
  }
}
