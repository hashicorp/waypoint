import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { GetLogStreamRequest } from 'waypoint-pb';
import { AppRouteModel } from '../app';

export default class Logs extends Route {
  @service api!: ApiService;

  async model() {
    const app = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    const latestDeployment = (await app.deployments).firstObject;
    const req = new GetLogStreamRequest();

    req.setDeploymentId(latestDeployment?.id!);

    return req;
  }
}
