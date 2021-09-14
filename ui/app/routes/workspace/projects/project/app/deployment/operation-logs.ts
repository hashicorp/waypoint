import Route from '@ember/routing/route';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';

export default class extends Route {
  breadcrumbs(): Breadcrumb[] {
    return [
      {
        label: 'Operation Logs',
        route: 'workspace.projects.project.app.deployment.operation-logs',
      },
    ];
  }
}
