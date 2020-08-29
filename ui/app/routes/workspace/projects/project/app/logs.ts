import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import CurrentApplicationService from 'waypoint/services/current-application';
import CurrentWorkspaceService from 'waypoint/services/current-workspace';

export default class Logs extends Route {
  @service api!: ApiService;
  @service currentApplication!: CurrentApplicationService;
  @service currentWorkspace!: CurrentWorkspaceService;

  breadcrumbs = [
    {
      label: 'Logs',
      args: ['workspace.projects.project.app.logs'],
    },
  ];

  async model() {
    // todo(pearkes): construct GetLogStreamRequest
  }
}
