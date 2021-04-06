import Route from '@ember/routing/route';
import { AppRouteModel } from '../app';
import { Project } from 'waypoint-pb';
export default class WorkspaceProjectsProjectAppSettings extends Route {
  breadcrumbs(model: AppRouteModel) {
    if (!model) return [];
    return [
      {
        label: model.application,
        icon: 'git-repository',
        args: ['workspace.projects.project.app'],
      },
      {
        label: 'Settings',
        icon: 'settings',
        args: ['workspace.projects.project.app.settings'],
      },
    ];
  }

  async model() {
    let app = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    return app.application;
  }

  setupController(controller: any, model: any, transition: any) {
    super.setupController(controller, model, transition);
    let proj = this.modelFor('workspace.projects.project') as Project.AsObject;
    controller.project = proj;
  }
}
