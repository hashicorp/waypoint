import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Ref, GetProjectRequest } from 'waypoint-pb';
interface ProjectModelParams {
  project_id: string;
}

export default class WorkspaceProjectsProjectSettings extends Route {
  @service api!: ApiService;

  breadcrumbs(model: Project) {
    if (!model) return [];
    return [
      {
        label: model.name,
        icon: 'folder-outline',
        args: ['workspace.projects.project.index'],
      },
      {
        label: 'Settings',
        icon: 'settings',
        args: ['workspace.projects.project.settings'],
      },
    ];
  }

  async model() {
    // Setup the project request
    let ref = new Ref.Project();
    let params = this.paramsFor('workspace.projects.project') as ProjectModelParams;
    ref.setProject(params.project_id);
    let req = new GetProjectRequest();
    req.setProject(ref);

    let resp = await this.api.client.getProject(req, this.api.WithMeta());
    let project = resp.getProject();

    return project?.toObject();
  }
}
