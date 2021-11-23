import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService, { DeploymentExtended } from 'waypoint/services/api';
import { Model as AppRouteModel } from '../app';

type Model = AppRouteModel['deployments'];

export default class DeploymentsList extends Route {
  @service api!: ApiService;

  async model(): Promise<DeploymentExtended[]> {
    let app = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    return app.deployments;
  }

  redirect(model: Model): void {
    if (model) {
      if (model[0]) {
        this.transitionTo('workspace.projects.project.app.deployment.deployment-seq', model[0].sequence);
      }
      this.transitionTo('workspace.projects.project.app.deployment');
    }
  }
}
