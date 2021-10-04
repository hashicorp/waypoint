import Route from '@ember/routing/route';
import { StatusReport } from 'waypoint-pb';
import { Model as ReleaseRouteModel } from '../release';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';
import { action } from '@ember/object';

interface Params {
  resource_id: string;
}

type Model = StatusReport.Resource.AsObject;

export default class extends Route {
  @action
  breadcrumbs(): Breadcrumb[] {
    let release = this.modelFor('workspace.projects.project.app.release') as ReleaseRouteModel;
    return [
      {
        label: `v${release.sequence}`,
        route: 'workspace.projects.project.app.release',
        icon: 'public-default',
      },
    ];
  }

  model({ resource_id }: Params): Model {
    let release = this.modelFor('workspace.projects.project.app.release') as ReleaseRouteModel;
    let resources = release.statusReport?.resourcesList ?? [];
    let resource = resources.find((r) => r.id === resource_id);

    if (!resource) {
      throw new Error(`Resource ${resource_id} not found`);
    }

    return resource;
  }
}
