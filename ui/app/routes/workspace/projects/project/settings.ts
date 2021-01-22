import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Project} from 'waypoint-pb';

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
    let proj = this.modelFor('workspace.projects.project') as Project.AsObject;
    return proj;
  }
}
