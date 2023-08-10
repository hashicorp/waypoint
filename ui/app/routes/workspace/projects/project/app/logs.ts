/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';
import { GetLogStreamRequest, Ref } from 'waypoint-pb';
import { Model as AppRouteModel } from '../app';
import { Model as WorkspaceRouteModel } from 'waypoint/routes/workspace';

type Model = {
  app: string;
  request: GetLogStreamRequest;
};

export default class Logs extends Route {
  @service api!: ApiService;

  breadcrumbs(model: Model): Breadcrumb[] {
    if (!model) return [];
    return [
      {
        label: model.app ?? 'unknown',
        route: 'workspace.projects.project.app',
      },
      {
        label: 'Logs',
        route: 'workspace.projects.project.app.logs',
      },
    ];
  }

  async model(): Promise<Model> {
    let app = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let ws = this.modelFor('workspace') as WorkspaceRouteModel;
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

    return { app: app.application.application, request: req };
  }
}
