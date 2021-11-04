import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Model as AppRouteModel } from '../app';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';
import { ReleaseExtended } from 'waypoint/services/api';

type Params = { sequence: string };
export type Model = ReleaseExtended;

export default class ReleaseDetail extends Route {
  @service api!: ApiService;

  breadcrumbs(model: Model): Breadcrumb[] {
    if (!model) return [];

    return [
      {
        label: model.application?.application ?? 'unknown',
        route: 'workspace.projects.project.app',
      },
      {
        label: 'Releases',
        route: 'workspace.projects.project.app.releases',
      },
    ];
  }

  async model(params: Params): Promise<Model> {
    let { releases } = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let release = releases.find((obj) => obj.sequence === Number(params.sequence));

    if (!release) {
      throw new Error(`Release v${params.sequence} not found`);
    }

    return release;
  }
}
