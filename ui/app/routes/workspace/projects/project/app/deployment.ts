import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Deployment, StatusReport, UI } from 'waypoint-pb';
import { Model as AppRouteModel } from '../app';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';

type Params = { sequence: string };
type Model = UI.DeploymentBundle.AsObject;

interface WithStatusReport {
  statusReport?: StatusReport.AsObject;
}

export default class DeploymentDetail extends Route {
  @service api!: ApiService;

  breadcrumbs(model: Model): Breadcrumb[] {
    if (!model) return [];
    return [
      {
        label: model.deployment?.application?.application ?? 'unknown',
        icon: 'git-repository',
        route: 'workspace.projects.project.app',
      },
      {
        label: 'Deployments',
        icon: 'upload',
        route: 'workspace.projects.project.app.deployments',
      },
      {
        label: `v${model.deployment?.sequence ?? 'unknown'}`,
        icon: 'upload',
        route: 'workspace.projects.project.app.deployment',
      },
    ];
  }

  async model(params: Params): Promise<Model> {
    let { deployments } = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let bundle = deployments.find((d) => d.deployment?.sequence === Number(params.sequence));

    if (!bundle) {
      throw new Error(`Deployment ${params.sequence} not found`);
    }

    return bundle;
  }

  afterModel(model: Deployment.AsObject & WithStatusReport): void {
    let { statusReports } = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let statusReport = statusReports.find((sr) => sr.deploymentId === model.id);

    model.statusReport = statusReport;
  }
}
