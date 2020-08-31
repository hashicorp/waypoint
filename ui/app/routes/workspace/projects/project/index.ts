import Route from '@ember/routing/route';
import { Project } from 'waypoint-pb';

export default class ProjectIndex extends Route {
  redirect(model: Project.AsObject) {
    return this.transitionTo('workspace.projects.project.apps', model.name);
  }
}
