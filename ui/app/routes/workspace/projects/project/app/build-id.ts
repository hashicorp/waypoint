import { Ref, GetBuildRequest } from 'waypoint-pb';
import BuildDetail from './build';

interface BuildModelIdParams {
  build_id: string;
}

export default class WorkspaceProjectsProjectAppBuildId extends BuildDetail {
  renderTemplate() {
    this.render('workspace/projects/project/app/build', {
      into: 'workspace/projects/project',
    });
  }

  async model(params: BuildModelIdParams) {
    // Setup the build request
    let ref = new Ref.Operation();
    ref.setId(params.build_id);
    let req = new GetBuildRequest();
    req.setRef(ref);

    let build = await this.api.client.getBuild(req, this.api.WithMeta());
    return build.toObject();
  }
}
