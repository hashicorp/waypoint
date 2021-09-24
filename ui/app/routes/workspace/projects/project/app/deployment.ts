import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { GetDeploymentRequest, Deployment, Ref, StatusReport } from 'waypoint-pb';
import { Model as AppRouteModel } from '../app';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';

type Params = { sequence: string };
export type Model = Deployment.AsObject & WithStatusReport;

interface WithStatusReport {
  statusReport?: StatusReport.AsObject;
}

export default class DeploymentDetail extends Route {
  @service api!: ApiService;

  breadcrumbs(model: Model): Breadcrumb[] {
    if (!model) return [];
    return [
      {
        label: model.application?.application ?? 'unknown',
        icon: 'git-repository',
        route: 'workspace.projects.project.app',
      },
    ];
  }

  async model(params: Params): Promise<Model> {
    let { deployments } = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let deploymentFromAppRoute = deployments.find((obj) => obj.sequence == Number(params.sequence));

    if (!deploymentFromAppRoute) {
      throw new Error(`Deployment v${params.sequence} not found`);
    }

    let ref = new Ref.Operation();
    ref.setId(deploymentFromAppRoute.id);
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
