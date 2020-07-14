import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import {ListBuildsRequest, Ref, ListBuildsResponse} from 'waypoint-pb';

export default class Builds extends Route {
  @service api!: ApiService;

  async model() {
    // The application we'd typically read from the users session
    // and an application list API but are doing this for proof
    // of concept
    var app = new Ref.Application()
    app.setProject("foobar")
    app.setApplication("wp-gcr-deno-test")    
  
    var req = new ListBuildsRequest()
    req.setApplication(app)

    var resp = await this.api.client.listBuilds(req, {})
    let buildResp: ListBuildsResponse = resp;

    return buildResp.getBuildsList().map(b => b.toObject());
  }
}
