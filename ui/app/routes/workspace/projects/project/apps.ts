import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Project, Application } from 'waypoint-pb';

export default class Apps extends Route {
  @service api!: ApiService;

  async model() {
    let proj = this.modelFor('workspace.projects.project') as Project.AsObject;
    return proj.applicationsList;
  }

  afterModel(model: Application.AsObject[]) {
    if (model.length == 1) {
      return this.transitionTo('workspace.projects.project.app', model[0].name);
    }
  }
}
