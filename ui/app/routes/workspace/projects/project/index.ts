import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { Project } from 'waypoint-pb';
import CurrentProjectService from 'waypoint/services/current-project';
import CurrentApplicationService from 'waypoint/services/current-application';

export default class ProjectIndex extends Route {
  @service currentProject!: CurrentProjectService;
  @service currentApplication!: CurrentApplicationService;

  redirect(model: Project.AsObject) {
    return this.transitionTo('workspace.projects.project.apps', model.name);
  }
}
