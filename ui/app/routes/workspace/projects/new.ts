import Route from '@ember/routing/route';
import { Project } from 'waypoint-pb';
export default class WorkspaceProjectsNew extends Route {
  breadcrumbs = [
    {
      label: 'Projects',
      args: ['workspace.projects'],
    },
    {
      label: 'New Project',
      args: ['workspace.projects.new'],
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
