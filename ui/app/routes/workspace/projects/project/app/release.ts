import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { GetReleaseRequest, Release, Ref, StatusReport } from 'waypoint-pb';
import { Model as AppRouteModel } from '../app';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';

type Params = { sequence: string };
export type Model = Release.AsObject & WithStatusReport;

interface WithStatusReport {
  statusReport?: StatusReport.AsObject;
}

export default class ReleaseDetail extends Route {
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
        label: 'Releases',
        icon: 'public-default',
        route: 'workspace.projects.project.app.releases',
      },
    ];
  }

  async model(params: Params): Promise<Model> {
    let { releases } = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let releaseFromAppRoute = releases.find((obj) => obj.sequence === Number(params.sequence));

    if (!releaseFromAppRoute) {
      throw new Error(`Release v${params.sequence} not found`);
    }

    let ref = new Ref.Operation();
    ref.setId(releaseFromAppRoute.id);
    let req = new GetReleaseRequest();
    req.setRef(ref);

    let release: Release = await this.api.client.getRelease(req, this.api.WithMeta());
    return release.toObject();
  }

  afterModel(model: Release.AsObject & WithStatusReport): void {
    let { statusReports } = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let statusReport = statusReports.find((sr) => sr.releaseId === model.id);

    model.statusReport = statusReport;
  }
}
