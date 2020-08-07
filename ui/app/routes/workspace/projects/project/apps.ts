import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import CurrentProjectService from 'waypoint/services/current-project';

export default class App extends Route {
  @service api!: ApiService;
  @service currentProject!: CurrentProjectService;

  async model() {
    return this.currentProject.project?.getApplicationsList().map((b) => b.toObject());
  }
}
