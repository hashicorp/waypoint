import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Ref, GetProjectRequest } from 'waypoint-pb';
import PollModelService from 'waypoint/services/poll-model';
interface ProjectModelParams {
  project_id: string;
}

export default class ProjectDetail extends Route {
  @service api!: ApiService;
  @service pollModel!: PollModelService;

  breadcrumbs = [
    {
      label: 'Projects',
      args: ['workspace.projects'],
    },
  ];

  async model(params: ProjectModelParams) {
    // Setup the project request
    let ref = new Ref.Project();
    ref.setProject(params.project_id);
    let req = new GetProjectRequest();
    req.setProject(ref);

    let resp = await this.api.client.getProject(req, this.api.WithMeta());
    let project = resp.getProject();

    return project?.toObject();
  }

  afterModel() {
    this.pollModel.setup(this);
  }
}
