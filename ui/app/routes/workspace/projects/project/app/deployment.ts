import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { GetDeploymentRequest, Deployment, Ref } from 'waypoint-pb';
import { AppRouteModel } from '../app';

interface DeploymentModelParams {
  deployment_id: string;
}

export default class DeploymentDetail extends Route {
  @service api!: ApiService;

  breadcrumbs(model: AppRouteModel) {
    if (!model) return [];
    return [
      {
        label: model.application.application,
        icon: 'git-repository',
        args: ['workspace.projects.project.app'],
      },
      {
        label: 'Deployments',
        args: ['workspace.projects.project.app.deployments'],
      },
    ];
  }

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
