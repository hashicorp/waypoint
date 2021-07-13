import Route from '@ember/routing/route';
import { Model as AppRouteModel } from '../app';
import DeploymentsController from 'waypoint/controllers/workspace/projects/project/app/deployments';

export default class Deployments extends Route {
  async model() {
    let app = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    return app.deployments;
  }

  resetController(controller: DeploymentsController, isExiting: boolean) {
    if (isExiting) {
      controller.set('isShowingDestroyed', null);
    }
  }
}
