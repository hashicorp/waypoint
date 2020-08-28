import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { GetReleaseRequest, Release, Ref } from 'waypoint-pb';
import CurrentApplicationService from 'waypoint/services/current-application';
import CurrentWorkspaceService from 'waypoint/services/current-workspace';

interface ReleaseModelParams {
  release_id: string;
}

export default class ReleaseDetail extends Route {
  @service api!: ApiService;
  @service currentApplication!: CurrentApplicationService;
  @service currentWorkspace!: CurrentWorkspaceService;

  async model(params: ReleaseModelParams) {
    var ref = new Ref.Operation();
    ref.setId(params.release_id);
    var req = new GetReleaseRequest();
    req.setRef(ref);

    var resp = await this.api.client.getRelease(req, {});
    let deploy: Release = resp;
    return deploy.toObject();
  }
}
