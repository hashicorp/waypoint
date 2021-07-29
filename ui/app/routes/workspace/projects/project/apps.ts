import Route from '@ember/routing/route';
import { Model as ProjectRouteModel } from 'waypoint/routes/workspace/projects/project';

type Model = ProjectRouteModel;

export default class Apps extends Route {
  async model(): Promise<Model> {
    // Technically we get this behavior for free because, in Ember, if a
    // route doesn’t define a model hook then it automatically receives
    // the model of its parent. Still, it doesn’t hurt to be explicit.
    return this.modelFor('workspace.projects.project') as ProjectRouteModel;
  }
}
