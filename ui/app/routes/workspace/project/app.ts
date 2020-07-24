import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Ref } from 'waypoint-pb';

interface AppModelParams {
  app_id: string;
}

export default class App extends Route {
  @service api!: ApiService;

  async model(params: AppModelParams) {
    let app = new Ref.Application();
    let proj = this.modelFor('workspace.project').ref;

    // App based on id
    app.setApplication(params.app_id);
    app.setProject(proj);

    return {
      ref: app as Ref.Application,
      app: app.toObject(),
    };
  }
}
