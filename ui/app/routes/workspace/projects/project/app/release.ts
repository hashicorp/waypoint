import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { GetReleaseRequest, Release, Ref, StatusReport } from 'waypoint-pb';
import { Model as AppRouteModel, Params as AppRouteParams } from '../app';
import { Params as ProjectRouteParams } from '../../project';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';

type Params = { sequence: string } | { release_id: string };
type Model = Release.AsObject;

interface WithStatusReport {
  statusReport?: StatusReport.AsObject;
}

export default class ReleaseDetail extends Route {
  @service api!: ApiService;

  breadcrumbs(model: Model): Breadcrumb[] {
    if (!model) return [];

    return [
      {
        label: model.application?.application ?? '',
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
    let req = new GetReleaseRequest();
    let ref = this.refForParams(params);

    req.setRef(ref);

    let release: Release = await this.api.client.getRelease(req, this.api.WithMeta());

    return release.toObject();
  }

  afterModel(model: Release.AsObject & WithStatusReport): void {
    let { statusReports } = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let statusReport = statusReports.find((sr) => sr.releaseId === model.id);

    model.statusReport = statusReport;
  }

  private refForParams(params: Params): Ref.Operation {
    if ('sequence' in params) {
      return this.refForSequence(params.sequence);
    } else {
      return this.refForId(params.release_id);
    }
  }

  private refForSequence(sequence: string): Ref.Operation {
    let ref = new Ref.Operation();
    let appRef = new Ref.Application();
    let seqRef = new Ref.OperationSeq();
    let { project_id } = this.paramsFor('workspace.projects.project') as ProjectRouteParams;
    let { app_id } = this.paramsFor('workspace.projects.project.app') as AppRouteParams;

    appRef.setProject(project_id);
    appRef.setApplication(app_id);

    seqRef.setApplication(appRef);
    seqRef.setNumber(Number(sequence));

    ref.setSequence(seqRef);

    return ref;
  }

  private refForId(id: string): Ref.Operation {
    let ref = new Ref.Operation();

    ref.setId(id);

    return ref;
  }
}
