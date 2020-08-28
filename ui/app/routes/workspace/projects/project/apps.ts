import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import CurrentProjectService from 'waypoint/services/current-project';
import { Project } from 'waypoint-pb';

export default class Apps extends Route {
  @service api!: ApiService;
  @service currentProject!: CurrentProjectService;

  async model() {
    return this.currentProject.project?.toObject();
  }

  afterModel(model: Project.AsObject) {
    if (model.applicationsList.length == 1) {
      return this.transitionTo('workspace.projects.project.app', model.applicationsList[0].name);
    }
  }
}
