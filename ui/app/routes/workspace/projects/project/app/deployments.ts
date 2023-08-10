/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { Model as AppRouteModel } from '../app';
import DeploymentsController from 'waypoint/controllers/workspace/projects/project/app/deployments';

export type Model = AppRouteModel['deployments'];

export default class Deployments extends Route {
  async model(): Promise<Model> {
    let app = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    return app.deployments;
  }

  resetController(controller: DeploymentsController, isExiting: boolean): void {
    if (isExiting) {
      controller.set('isShowingDestroyed', null);
    }
  }
}
