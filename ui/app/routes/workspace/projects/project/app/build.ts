import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Ref, GetBuildRequest } from 'waypoint-pb';
import { Model as AppRouteModel } from '../app';

interface BuildModelParams {
  sequence: number;
}
export default class BuildDetail extends Route {
  @service api!: ApiService;

  breadcrumbs(model: AppRouteModel) {
    if (!model) return [];
    return [
      {
        label: model.application.application,
        icon: 'git-repository',
        args: ['workspace.projects.project.app'],
      },
      {
        label: 'Builds',
        icon: 'build',
        args: ['workspace.projects.project.app.builds'],
      },
    ];
  }

  async model(params: BuildModelParams) {
    // Setup the build request
    let { builds } = this.modelFor('workspace.projects.project.app');
    let { id: build_id } = builds.find((obj) => obj.sequence === Number(params.sequence));

    let ref = new Ref.Operation();
    ref.setId(build_id);
    let req = new GetBuildRequest();
    req.setRef(ref);

    let build = await this.api.client.getBuild(req, this.api.WithMeta());
    return build.toObject();
  }
}
