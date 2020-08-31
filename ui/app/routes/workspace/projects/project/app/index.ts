import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { Ref } from 'waypoint-pb';
import CurrentProjectService from 'waypoint/services/current-project';
import CurrentApplicationService from 'waypoint/services/current-application';

export default class AppIndex extends Route {
  @service currentProject!: CurrentProjectService;
  @service currentApplication!: CurrentApplicationService;

  redirect(model: Ref.Application.AsObject) {
    return this.transitionTo('workspace.projects.project.app.logs', model.application);
  }
}
