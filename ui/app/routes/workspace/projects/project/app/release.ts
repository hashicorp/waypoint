import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { GetReleaseRequest, Release, Ref } from 'waypoint-pb';

interface ReleaseModelParams {
  release_id: string;
}

export default class ReleaseDetail extends Route {
  @service api!: ApiService;

  async model(params: ReleaseModelParams) {
    var ref = new Ref.Operation();
    ref.setId(params.release_id);
    var req = new GetReleaseRequest();
    req.setRef(ref);

    var resp = await this.api.client.getRelease(req, this.api.WithMeta());
    let deploy: Release = resp;
    return deploy.toObject();
  }
}
