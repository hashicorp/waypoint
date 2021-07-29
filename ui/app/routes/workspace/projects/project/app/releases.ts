import Route from '@ember/routing/route';
import { Model as AppRouteModel } from '../app';

type Model = AppRouteModel['releases'];

export default class Releases extends Route {
  async model(): Promise<Model> {
    let app = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    return app.releases;
  }
}
