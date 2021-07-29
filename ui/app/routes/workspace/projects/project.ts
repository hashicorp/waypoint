import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { GetProjectRequest, Project, Ref } from 'waypoint-pb';
import PollModelService from 'waypoint/services/poll-model';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';

export type Params = { project_id: string };
export type Model = Project.AsObject;

export default class ProjectDetail extends Route {
  @service api!: ApiService;
  @service pollModel!: PollModelService;

  breadcrumbs: Breadcrumb[] = [
    {
      label: 'Projects',
      route: 'workspace.projects',
    },
  ];

  async model(params: Params): Promise<Model> {
    // Setup the project request
    let ref = new Ref.Project();
    ref.setProject(params.project_id);
    let req = new GetProjectRequest();
    req.setProject(ref);

    let resp = await this.api.client.getProject(req, this.api.WithMeta());
    let project = resp.getProject();

    if (!project) {
      // In reality the API will return an error in this circumstance
      // but the types don’t tell us that.
      throw new Error(`Project ${params.project_id} not found`);
    }

    return project.toObject();
  }

  afterModel(): void {
    this.pollModel.setup(this);
  }
}
