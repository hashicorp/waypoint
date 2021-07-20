import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { GetReleaseRequest, Release, Ref, StatusReport } from 'waypoint-pb';
import { Model as AppRouteModel } from '../app';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';

interface ReleaseModelParams {
  sequence: number;
}

interface WithStatusReport {
  statusReport?: StatusReport.AsObject;
}

export default class ReleaseDetail extends Route {
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
        label: 'Releases',
        icon: 'public-default',
        route: 'workspace.projects.project.app.releases',
      },
    ];
  }

  async model(params: ReleaseModelParams): Promise<Release.AsObject> {
    let { releases } = this.modelFor('workspace.projects.project.app');
    let { id: release_id } = releases.find((obj) => obj.sequence === Number(params.sequence));

    let ref = new Ref.Operation();
    ref.setId(release_id);
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
