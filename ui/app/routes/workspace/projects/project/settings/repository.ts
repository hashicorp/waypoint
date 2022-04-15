import Route from '@ember/routing/route';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';

export default class WorkspaceProjectsProjectSettingsRepository extends Route {
  breadcrumbs(): Breadcrumb[] {
    return [
      {
        label: 'Git Repository',
        route: 'workspace.projects.project.settings.config-variables',
      },
    ];
  }
}
