import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import CurrentApplicationService from 'waypoint/services/current-application';
import CurrentWorkspaceService from 'waypoint/services/current-workspace';
import { ListReleasesRequest, ListReleasesResponse } from 'waypoint-pb';

export default class Releases extends Route {
  @service api!: ApiService;
  @service currentApplication!: CurrentApplicationService;
  @service currentWorkspace!: CurrentWorkspaceService;

  async model() {
    var req = new ListReleasesRequest();
    req.setApplication(this.currentApplication.ref);
    req.setWorkspace(this.currentWorkspace.ref);

    var resp = await this.api.client.listReleases(req, this.api.WithMeta());
    let releaseResp: ListReleasesResponse = resp;

    return releaseResp.getReleasesList().map((r) => r.toObject());
  }
}
