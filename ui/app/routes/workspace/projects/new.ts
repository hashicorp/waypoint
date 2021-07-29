import Route from '@ember/routing/route';
import { Project } from 'waypoint-pb';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';
import Controller from 'waypoint/controllers/workspace/projects/new';

type Model = Project.AsObject;

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

  model(): Model {
    let proj = new Project();
    return proj.toObject();
  }

  resetController(controller: Controller, isExiting: boolean): void {
    if (isExiting) {
      controller.set('createGit', false);
    }
  }
}
