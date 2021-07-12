import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Ref, Deployment, Build, Release, Project, StatusReport } from 'waypoint-pb';
import PollModelService from 'waypoint/services/poll-model';
import ObjectProxy from '@ember/object/proxy';
import PromiseProxyMixin from '@ember/object/promise-proxy-mixin';
import { resolve, hash } from 'rsvp';

interface AppModelParams {
  app_id: string;
}

export interface AppRouteModel {
  application: Ref.Application.AsObject;
  deployments: Promise<Deployment.AsObject[]>;
  releases: Promise<Release.AsObject[]>;
  builds: Promise<Build.AsObject[]>;
  statusReports: Promise<StatusReport.AsObject[]>;
}

interface WithStatusReport {
  statusReport: StatusReport.AsObject;
}

export interface ResolvedModel {
  application: Ref.Application.AsObject;
  deployments: (Deployment.AsObject & WithStatusReport)[];
  releases: (Release.AsObject & WithStatusReport)[];
  builds: (Build.AsObject & WithStatusReport)[];
  statusReports: StatusReport.AsObject[];
}

export default class App extends Route {
  @service api!: ApiService;
  @service pollModel!: PollModelService;

  breadcrumbs(model: AppRouteModel) {
    if (!model) return [];
    return [
      {
        label: model.application.project!,
        icon: 'folder-outline',
        args: ['workspace.projects.project.apps'],
      },
    ];
  }

  async model(params: AppModelParams): Promise<AppRouteModel> {
    let ws = this.modelFor('workspace') as Ref.Workspace.AsObject;
    let wsRef = new Ref.Workspace();
    wsRef.setWorkspace(ws.workspace);

    let proj = this.modelFor('workspace.projects.project') as Project.AsObject;

    let appRef = new Ref.Application();
    // App based on id
    appRef.setApplication(params.app_id);
    appRef.setProject(proj.name);

    let ObjectPromiseProxy = ObjectProxy.extend(PromiseProxyMixin);

    return hash({
      application: appRef.toObject(),
      deployments: ObjectPromiseProxy.create({
        promise: resolve(this.api.listDeployments(wsRef, appRef)),
      }),
      releases: ObjectPromiseProxy.create({
        promise: resolve(this.api.listReleases(wsRef, appRef)),
      }),
      builds: ObjectPromiseProxy.create({
        promise: resolve(this.api.listBuilds(wsRef, appRef)),
      }),
      statusReports: ObjectPromiseProxy.create({
        promise: resolve(this.api.listStatusReports(wsRef, appRef)),
      }),
      latestStatusReport: ObjectPromiseProxy.create({
        promise: resolve(this.api.getLatestStatusReport(wsRef, appRef)),
      }),
    });
  }

  afterModel(model: ResolvedModel): void {
    injectStatusReports(model);
    // this.pollModel.setup(this);
  }
}

function injectStatusReports(model: ResolvedModel): void {
  let { deployments, releases, statusReports } = model;

  for (let statusReport of statusReports) {
    if (statusReport.deploymentId) {
      let deployment = deployments.find((d) => d.id === statusReport.deploymentId);
      if (deployment) {
        deployment.statusReport = statusReport;
      }
    } else if (statusReport.releaseId) {
      let release = releases.find((d) => d.id === statusReport.releaseId);
      if (release) {
        release.statusReport = statusReport;
      }
    }
  }
}
