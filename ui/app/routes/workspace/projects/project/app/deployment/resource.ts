import Route from '@ember/routing/route';
import { StatusReport } from 'waypoint-pb';
import { Model as DeploymentRouteModel } from '../deployment';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';
import { action } from '@ember/object';

interface Params {
  resource_id: string;
}

type Model = StatusReport.Resource.AsObject;

export default class extends Route {
  @action
  breadcrumbs(): Breadcrumb[] {
    let deployment = this.modelFor('workspace.projects.project.app.deployment') as DeploymentRouteModel;
    return [
      {
        label: `v${deployment.sequence}`,
        route: 'workspace.projects.project.app.deployment',
        icon: 'upload',
      },
    ];
  }

  model({ resource_id }: Params): Model {
    let deployment = this.modelFor('workspace.projects.project.app.deployment') as DeploymentRouteModel;
    let resources = deployment.statusReport?.resourcesList ?? [];
    let resource = resources.find((r) => r.id === resource_id);

    if (!resource) {
      throw new Error(`Resource ${resource_id} not found`);
    }

    return resource;
  }
}
