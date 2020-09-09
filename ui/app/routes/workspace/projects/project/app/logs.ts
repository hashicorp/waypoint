import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';

export default class Logs extends Route {
  @service api!: ApiService;

  async model() {
    // todo(pearkes): construct GetLogStreamRequest
  }
}
