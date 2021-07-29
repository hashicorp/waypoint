import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';

export default class Workspaces extends Route {
  @service api!: ApiService;

  redirect(): void {
    // For now, we just support the default workspace
    this.transitionTo('workspace', 'default');
  }
}
