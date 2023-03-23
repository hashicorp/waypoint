/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApiService from 'waypoint/services/api';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';
import { Model as AppRouteModel } from '../app';
import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class Exec extends Route {
  @service api!: ApiService;

  breadcrumbs(model: AppRouteModel): Breadcrumb[] {
    if (!model) return [];
    return [
      {
        label: model.application.application ?? 'unknown',
        route: 'workspace.projects.project.app',
      },
      {
        label: 'Exec',
        route: 'workspace.projects.project.app.exec',
      },
    ];
  }

  async model(): Promise<AppRouteModel> {
    let app = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    return app;
  }
}
