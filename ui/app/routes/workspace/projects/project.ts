import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { UI, Project, Ref, Job } from 'waypoint-pb';
import PollModelService from 'waypoint/services/poll-model';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';
import ProjectService from 'waypoint/services/project';

export type Params = { project_id: string };
export type Model = Project.AsObject & {
  latestInitJob?: Job.AsObject;
};

export default class ProjectDetail extends Route {
  @service api!: ApiService;
  @service pollModel!: PollModelService;
  @service declare project: ProjectService;

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
    let req = new UI.GetProjectRequest();
    req.setProject(ref);

    let resp = await this.api.client.uI_GetProject(req, this.api.WithMeta());
    let project = resp.getProject();

    if (!project) {
      // In reality the API will return an error in this circumstance
      // but the types don’t tell us that.
      throw new Error(`Project ${params.project_id} not found`);
    }

    let result = project.toObject() as Model;

    // TODO(jgwhite): It’d be better not to sneak this onto the project,
    // but changing the model type of this route is way more disruptive.
    result.latestInitJob = resp.getLatestInitJob()?.toObject();

    return result;
  }

  afterModel(model: Model): void {
    this.pollModel.setup(this);
    this.project.current = model;
  }

  deactivate() {
    this.project.current = undefined;
  }

  willDestroy() {
    this.project.current = undefined;
  }
}
