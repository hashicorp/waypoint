import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';

export default class Exec extends Route {
  @service api!: ApiService;

  async model(): Promise<void> {
    // todo(pearkes): construct GetExecStreamRequest
  }
}
