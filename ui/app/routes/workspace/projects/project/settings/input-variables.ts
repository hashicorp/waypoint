import Route from '@ember/routing/route';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';

export default class WorkspaceProjectsProjectSettingsInputVariables extends Route {
  breadcrumbs(): Breadcrumb[] {
    return [
      {
        label: 'Input Variables',
        route: 'workspace.projects.project.settings.config-variables',
      },
    ];
  }
}
