import Transition from '@ember/routing/-private/transition';
import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Ref } from 'waypoint-pb';
import CurrentProjectService from 'waypoint/services/current-project';
import CurrentApplicationService from 'waypoint/services/current-application';

interface AppModelParams {
  app_id: string;
}

export default class App extends Route {
  @service api!: ApiService;
  @service currentProject!: CurrentProjectService;
  @service currentApplication!: CurrentApplicationService;

  async model(params: AppModelParams) {
    let app = new Ref.Application();
    let proj = this.currentProject.ref;

    // App based on id
    app.setApplication(params.app_id);
    app.setProject(proj?.getProject()!);

    // Set ref on current app
    this.currentApplication.ref = app;

    return app.toObject();
  }

  afterModel(model: Ref.Application.AsObject, transition: Transition) {
    if (model) {
      return this.transitionTo('workspace.projects.project.app.logs', model.application);
    }
    return super.afterModel(model, transition);
  }
}
