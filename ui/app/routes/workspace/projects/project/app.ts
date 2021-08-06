import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Ref, Deployment, Build, Release, Project, StatusReport, PushedArtifact } from 'waypoint-pb';
import PollModelService from 'waypoint/services/poll-model';
import { hash } from 'rsvp';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';

export interface Params {
  app_id: string;
}

export interface Model {
  application: Ref.Application.AsObject;
  deployments: (Deployment.AsObject & WithStatusReport)[];
  releases: (Release.AsObject & WithStatusReport)[];
  builds: (Build.AsObject & WithPushedArtifact)[];
  pushedArtifacts: PushedArtifact.AsObject[];
  statusReports: StatusReport.AsObject[];
}

interface WithStatusReport {
  statusReport?: StatusReport.AsObject;
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
      application: appRef.toObject(),
      deployments: this.api.listDeployments(wsRef, appRef),
      releases: this.api.listReleases(wsRef, appRef),
      builds: this.api.listBuilds(wsRef, appRef),
      pushedArtifacts: this.api.listPushedArtifacts(wsRef, appRef),
      statusReports: this.api.listStatusReports(wsRef, appRef),
    });
  }

  afterModel(model: Model): void {
    injectStatusReports(model);
    injectPushedArtifacts(model);
    this.pollModel.setup(this);
  }
}

function injectStatusReports(model: Model): void {
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
