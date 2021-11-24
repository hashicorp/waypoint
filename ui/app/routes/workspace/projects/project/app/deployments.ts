import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService, { DeploymentExtended } from 'waypoint/services/api';
import { Model as AppRouteModel } from '../app';

export default class DeploymentsList extends Route {
  @service api!: ApiService;

  async model(): Promise<DeploymentExtended[]> {
    let app = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    return app.deployments;
  }

  redirect(): void {
    this.transitionTo('workspace.projects.project.app.deployment');
  }
}
