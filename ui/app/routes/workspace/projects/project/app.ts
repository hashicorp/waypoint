import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Ref, Deployment, Build, Release, Project } from 'waypoint-pb';
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
    });
  }

  afterModel() {
    this.pollModel.setup(this);
  }
}
