import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { GetLogStreamRequest, Ref } from 'waypoint-pb';
import { AppRouteModel } from '../app';

export default class Logs extends Route {
  @service api!: ApiService;

  async model() {
    const app = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let ws = this.modelFor('workspace') as Ref.Workspace.AsObject;

    const req = new GetLogStreamRequest();
    const appReq = new GetLogStreamRequest.Application();

    const appRef = new Ref.Application();
    appRef.setApplication(app.application.application);
    const wsRef = new Ref.Workspace();
    wsRef.setWorkspace(ws.workspace);

    appReq.setApplication(appRef);
    appReq.setWorkspace(wsRef);
    appReq.setWorkspace(wsRef);

    req.setApplication(appReq);

    return req;
  }
}
