import ApiService from 'waypoint/services/api';
import { Model as AppRouteModel } from '../app';
import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class Exec extends Route {
  @service api!: ApiService;

  async model(): Promise<AppRouteModel> {
    let app = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    return app;
    // todo(pearkes): construct GetExecStreamRequest
  }
}
