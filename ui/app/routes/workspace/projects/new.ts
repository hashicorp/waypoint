import Route from '@ember/routing/route';
import { Project } from 'waypoint-pb';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';

export default class WorkspaceProjectsNew extends Route {
  breadcrumbs: Breadcrumb[] = [
    {
      label: 'Projects',
      route: 'workspace.projects',
    },
    {
      label: 'New Project',
      route: 'workspace.projects.new',
    },
  ];

  model() {
    let proj = new Project();
    return proj.toObject();
  }

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.set('createGit', false);
    }
  }
}
