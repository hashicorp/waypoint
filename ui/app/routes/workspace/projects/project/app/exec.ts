import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { ExecStreamRequest } from 'waypoint-pb';

export default class Exec extends Route {
  @service api!: ApiService;

  async model() {
    let application = this.modelFor('workspace.projects.project.app');
    let latestDeploymentId = application.releases[0]?.deploymentId;
    let req = new ExecStreamRequest();
    let start = new ExecStreamRequest.Start();
    debugger;
    start.setDeploymentId(latestDeploymentId);
    // this.api.client.
    // todo(pearkes): construct GetExecStreamRequest
  }
}
