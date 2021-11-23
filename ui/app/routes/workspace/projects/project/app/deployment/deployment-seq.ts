import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Model as AppRouteModel } from '../../app';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';
import { DeploymentExtended, ReleaseExtended } from 'waypoint/services/api';

type Params = { sequence: string };
export type Model = DeploymentExtended & WithRelease;

interface WithRelease {
  release?: ReleaseExtended;
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
      {
        label: 'Deployments',
        icon: 'upload',
        route: 'workspace.projects.project.app.deployments',
      },
    ];
  }

  async model(params: Params): Promise<Model> {
    let { deployments } = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let deployment = deployments.find((obj) => obj.sequence == Number(params.sequence));

    if (!deployment) {
      throw new Error(`Deployment v${params.sequence} not found`);
    }

    let deploymentId = deployment.id;
    let { releases } = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let release = releases.find((r) => r.deploymentId === deploymentId);
    return { ...deployment, release };
  }
}
