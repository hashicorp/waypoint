import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService, { DeploymentExtended, ReleaseExtended } from 'waypoint/services/api';
import { Ref, Build, Project, PushedArtifact, Workspace } from 'waypoint-pb';
import PollModelService from 'waypoint/services/poll-model';
import { hash } from 'rsvp';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';

export interface Params {
  app_id: string;
}

export interface Model {
  project: Project.AsObject;
  application: Ref.Application.AsObject;
  deployments: DeploymentExtended[];
  releases: ReleaseExtended[];
  builds: (Build.AsObject & WithPushedArtifact)[];
  pushedArtifacts: PushedArtifact.AsObject[];
  workspace: Ref.Workspace.AsObject;
  workspaces: Workspace.AsObject[];
}

interface WithPushedArtifact {
  pushedArtifact?: PushedArtifact.AsObject;
}

export default class App extends Route {
  @service api!: ApiService;
  @service pollModel!: PollModelService;

  breadcrumbs(model: Model): Breadcrumb[] {
    if (!model) return [];

    return [
      {
        label: model.application.project,
        icon: 'folder-outline',
        route: 'workspace.projects.project.apps',
      },
    ];
  }

  async model(params: Params): Promise<Model> {
    let ws = this.modelFor('workspace') as Ref.Workspace.AsObject;
    let wsRef = new Ref.Workspace();
    wsRef.setWorkspace(ws.workspace);

    let proj = this.modelFor('workspace.projects.project') as Project.AsObject;

    let appRef = new Ref.Application();
    // App based on id
    appRef.setApplication(params.app_id);
    appRef.setProject(proj.name);

    return hash({
      project: proj,
      application: appRef.toObject(),
      deployments: this.api.listDeployments(wsRef, appRef),
      releases: this.api.listReleases(wsRef, appRef),
      builds: this.api.listBuilds(wsRef, appRef),
      pushedArtifacts: this.api.listPushedArtifacts(wsRef, appRef),
      workspace: ws,
      workspaces: this.api.listWorkspaces(),
    });
  }

  afterModel(model: Model): void {
    injectPushedArtifacts(model);
    this.pollModel.setup(this);
  }
}

function injectPushedArtifacts(model: Model): void {
  let { builds, pushedArtifacts } = model;

  for (let pushedArtifact of pushedArtifacts) {
    if (pushedArtifact.buildId) {
      let build = builds.find((b) => b.id === pushedArtifact.buildId);
      if (build) {
        build.pushedArtifact = pushedArtifact;
      }
    }
  }
}
