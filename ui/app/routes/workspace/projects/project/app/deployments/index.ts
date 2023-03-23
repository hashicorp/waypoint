/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Model } from '../deployments';

export default class DeploymentsList extends Route {
  @service api!: ApiService;

  redirect(model: Model): void {
    if (model.length === 0) {
      return;
    }

    let latestDeployment = model[0];
    this.transitionTo('workspace.projects.project.app.deployments.deployment', latestDeployment.sequence);
  }
}
