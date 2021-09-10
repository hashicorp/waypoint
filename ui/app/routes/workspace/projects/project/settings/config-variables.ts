import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { ConfigGetRequest, ConfigGetResponse, Ref } from 'waypoint-pb';
import { Params as ProjectRouteParams } from 'waypoint/routes/workspace/projects/project';
export default class WorkspaceProjectsProjectSettingsConfigVariables extends Route {
  @service api!: ApiService;

  async model(): Promise<ConfigGetResponse.AsObject> {
    let ref = new Ref.Project();
    let params = this.paramsFor('workspace.projects.project') as ProjectRouteParams;
    ref.setProject(params.project_id);
    let req = new ConfigGetRequest();
    req.setProject(ref);

    let config = await this.api.client.getConfig(req, this.api.WithMeta());
    return config?.toObject();
  }

  setupController(controller, model, transition) {
    super.setupController(controller, model, transition);
    let project = this.modelFor('workspace.projects.project');

    controller.project = project;
  }
}
