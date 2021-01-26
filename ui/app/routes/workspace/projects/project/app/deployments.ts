import Route from '@ember/routing/route';
import { AppRouteModel } from '../app';


interface DeploymentsModelParams {
  destroyed: boolean;
}
export default class Deployments extends Route {
  async model(params: DeploymentsModelParams) {
    let app = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    return app.deployments;
  }
}
