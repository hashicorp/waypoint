import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Ref, GetProjectRequest, Project } from 'waypoint-pb';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';

interface ProjectModelParams {
  project_id: string;
}

export default class WorkspaceProjectsProjectSettings extends Route {
  @service api!: ApiService;

  breadcrumbs(model: Project.AsObject): Breadcrumb[] {
    if (!model) return [];
    return [
      {
        label: model.name,
        icon: 'folder-outline',
        route: 'workspace.projects.project.index',
      },
      {
        label: 'Settings',
        icon: 'settings',
        route: 'workspace.projects.project.settings',
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
