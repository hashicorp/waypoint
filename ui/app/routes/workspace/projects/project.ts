import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Ref, GetProjectRequest } from 'waypoint-pb';
import CurrentProjectService from 'waypoint/services/current-project';

interface ProjectModelParams {
  project_id: string;
}

export default class Project extends Route {
  @service api!: ApiService;
  @service currentProject!: CurrentProjectService;

  async model(params: ProjectModelParams) {
    // Setup the project request
    let ref = new Ref.Project();
    ref.setProject(params.project_id);
    let req = new GetProjectRequest();
    req.setProject(ref);

    // Set on service
    this.currentProject.ref = ref;

    let resp = await this.api.client.getProject(req, this.api.WithMeta());
    let project = resp.getProject();

    // Set on service
    this.currentProject.project = project;

    return project?.toObject();
  }
}
