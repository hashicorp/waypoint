import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Ref, Deployment, Build, Release, Project } from 'waypoint-pb';

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

  breadcrumbs(model: AppRouteModel) {
    if (!model) return [];
    return [
      {
        label: model.application.project!,
        args: ['workspace.projects.project.apps'],
      },
      {
        label: 'Application',
        args: ['workspace.projects.project.app'],
      },
      {
        label: model.application.application!,
        args: ['workspace.projects.project.app'],
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

    return {
      application: appRef.toObject(),
      deployments: this.api.listDeployments(wsRef, appRef),
      releases: this.api.listReleases(wsRef, appRef),
      builds: this.api.listBuilds(wsRef, appRef),
    };
  }
}
