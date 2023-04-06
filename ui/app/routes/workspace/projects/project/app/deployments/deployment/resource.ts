/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { StatusReport } from 'waypoint-pb';
import { Model as DeploymentRouteModel } from '../deployment';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';
import { action } from '@ember/object';

interface Params {
  resource_id: string;
}

type Model = StatusReport.Resource.AsObject;

export default class extends Route {
  @action
  breadcrumbs(model: Model): Breadcrumb[] {
    return [
      {
        label: 'Resources',
        route: 'workspace.projects.project.app.deployments.deployment',
      },
      {
        label: model.name,
        route: 'workspace.projects.project.app.deployments.deployment.resource',
      },
    ];
  }

  model({ resource_id }: Params): Model {
    let deployment = this.modelFor(
      'workspace.projects.project.app.deployments.deployment'
    ) as DeploymentRouteModel;
    let resources = deployment.statusReport?.resourcesList ?? [];
    let resource = resources.find((r) => r.id === resource_id);

    if (!resource) {
      throw new Error(`Resource ${resource_id} not found`);
    }

    return resource;
  }
}
