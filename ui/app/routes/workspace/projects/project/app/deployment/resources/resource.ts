import Route from '@ember/routing/route';
import { StatusReport } from 'waypoint-pb';
import { Model as DeploymentRouteModel } from '../../deployment';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';

interface Params {
  resource_id: string;
}

type Model = StatusReport.Resource.AsObject;

export default class extends Route {
  breadcrumbs(model: Model): Breadcrumb[] {
    return [
      {
        label: model.name,
        route: 'workspace.projects.project.app.deployment.resources.resource',
      },
    ];
  }

  model({ resource_id }: Params): StatusReport.Resource.AsObject {
    let bundle = this.modelFor('workspace.projects.project.app.deployment') as DeploymentRouteModel;
    let resources = bundle.latestStatusReport?.resourcesList ?? [];
    let resource = resources.find((r) => r.id === resource_id);

    if (!resource) {
      throw new Error(`Resource ${resource_id} not found`);
    }

    return resource;
  }
}
