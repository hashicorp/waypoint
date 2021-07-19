import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { GetLogStreamRequest, Ref } from 'waypoint-pb';
import { Model as AppRouteModel } from '../app';

export default class Logs extends Route {
  @service api!: ApiService;

  async model() {
    let app = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let ws = this.modelFor('workspace') as Ref.Workspace.AsObject;
    let req = new GetLogStreamRequest();
    let appReq = new GetLogStreamRequest.Application();

    let appRef = new Ref.Application();
    appRef.setApplication(app.application.application);
    appRef.setProject(app.application.project);
    let wsRef = new Ref.Workspace();
    wsRef.setWorkspace(ws.workspace);

    appReq.setApplication(appRef);
    appReq.setWorkspace(wsRef);

    req.setApplication(appReq);

    return req;
  }
}
