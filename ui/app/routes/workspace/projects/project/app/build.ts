import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Ref, GetBuildRequest } from 'waypoint-pb';

interface BuildModelParams {
  build_id: string;
}
export default class BuildDetail extends Route {
  @service api!: ApiService;

  async model(params: BuildModelParams) {
    // Setup the build request
    let ref = new Ref.Operation();
    ref.setId(params.build_id);
    let req = new GetBuildRequest();
    req.setRef(ref);

    let build = await this.api.client.getBuild(req, this.api.WithMeta());

    return build.toObject();
  }
}
