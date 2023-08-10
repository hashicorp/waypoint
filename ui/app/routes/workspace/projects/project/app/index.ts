/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { Model as AppRouteModel } from '../app';

export default class AppIndex extends Route {
  redirect(model: AppRouteModel): void {
    this.transitionTo('workspace.projects.project.app.deployments', model.application.application);
  }
}
