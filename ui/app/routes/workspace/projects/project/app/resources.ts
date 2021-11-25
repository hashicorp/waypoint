import { Model as AppRouteModel } from '../app';
import Route from '@ember/routing/route';

type Model = AppRouteModel['deployments'];

export default class Resources extends Route {
  async model(): Promise<Model> {

    let app = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let deployments = app.deployments;
    let resources = app.deployments.flatMap(dep => dep.statusReport?.resourcesList);
    let model = { deployments, resources };
    return model;
  }
}
