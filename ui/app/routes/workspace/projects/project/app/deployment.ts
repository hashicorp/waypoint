import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { GetDeploymentRequest, Deployment, Ref, StatusReport } from 'waypoint-pb';
import { Model as AppRouteModel } from '../app';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';

interface DeploymentModelParams {
  sequence: number;
}

interface WithStatusReport {
  statusReport?: StatusReport.AsObject;
}

export default class DeploymentDetail extends Route {
  @service api!: ApiService;

  breadcrumbs(model: AppRouteModel): Breadcrumb[] {
    if (!model) return [];
    return [
      {
        label: model.application.application,
        icon: 'git-repository',
        route: 'workspace.projects.project.app',
      },
      {
        label: 'Deployments',
        icon: 'upload',
        route: 'workspace.projects.project.app.deployments',
      },
    ];
  }

  async model(params: DeploymentModelParams): Promise<Deployment.AsObject> {
    let { deployments } = this.modelFor('workspace.projects.project.app');
    let { id: deployment_id } = deployments.find((obj) => obj.sequence == Number(params.sequence));

    let ref = new Ref.Operation();
    ref.setId(deployment_id);
    let req = new GetDeploymentRequest();
    req.setRef(ref);

    let resp = await this.api.client.getDeployment(req, this.api.WithMeta());
    let deploy: Deployment = resp;
    return deploy.toObject();
  }

  afterModel(model: Deployment.AsObject & WithStatusReport): void {
    let { statusReports } = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let statusReport = statusReports.find((sr) => sr.deploymentId === model.id);

    model.statusReport = statusReport;
  }
}
