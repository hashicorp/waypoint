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
  breadcrumbs(model: Model): Breadcrumb[] {
    return [
      {
        label: 'Resources',
        route: 'workspace.projects.project.app.release',
      },
      {
        label: model.name,
        route: 'workspace.projects.project.app.release.resource',
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
