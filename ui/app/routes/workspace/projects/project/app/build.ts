import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Ref, GetBuildRequest, Build, PushedArtifact } from 'waypoint-pb';
import { Model as AppRouteModel } from '../app';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';

type Params = { sequence: string };
type Model = Build.AsObject & WithPushedArtifact;

interface WithPushedArtifact {
  pushedArtifact?: PushedArtifact.AsObject;
}

type BuildWithArtifact = Build.AsObject & WithPushedArtifact;

export default class BuildDetail extends Route {
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
        label: 'Builds',
        icon: 'build',
        route: 'workspace.projects.project.app.builds',
      },
    ];
  }

  async model(params: Params): Promise<Model> {
    let { builds } = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let buildFromAppRoute = builds.find((obj) => obj.sequence === Number(params.sequence));

    if (!buildFromAppRoute) {
      throw new Error(`Build v${params.sequence} not found`);
    }

    let ref = new Ref.Operation();
    ref.setId(buildFromAppRoute.id);
    let req = new GetBuildRequest();
    req.setRef(ref);

    let build = await this.api.client.getBuild(req, this.api.WithMeta());
    let result: BuildWithArtifact = build.toObject();

    result.pushedArtifact = buildFromAppRoute.pushedArtifact;

    return result;
  }
}
