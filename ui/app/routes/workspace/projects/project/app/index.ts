import Route from '@ember/routing/route';
import { Model as AppRouteModel } from '../app';

export default class AppIndex extends Route {
  redirect(model: AppRouteModel): void {
    this.transitionTo('workspace.projects.project.app.logs', model.application.application);
  }
}
