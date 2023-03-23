/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { Model as AppRouteModel } from '../app';

export default class AppIndex extends Route {
  redirect(model: AppRouteModel): void {
    this.transitionTo('workspace.projects.project.app.deployments', model.application.application);
  }
}
