import Route from '@ember/routing/route';
import { AppRouteModel } from '../app';

export default class Releases extends Route {
  async model() {
    let app = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    return app.releases;
  }
}
