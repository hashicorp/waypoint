/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';
import { Model as ProjectRouteModel } from 'waypoint/routes/workspace/projects/project';

type Model = ProjectRouteModel;

export default class Apps extends Route {
  breadcrumbs(model: Model): Breadcrumb[] {
    if (!model) return [];

    return [
      {
        label: model.name,
        route: 'workspace.projects.project.apps',
      },
      {
        label: 'Applications',
        route: 'workspace.projects.project.apps',
      },
    ];
  }

  async model(): Promise<Model> {
    // Technically we get this behavior for free because, in Ember, if a
    // route doesn’t define a model hook then it automatically receives
    // the model of its parent. Still, it doesn’t hurt to be explicit.
    return this.modelFor('workspace.projects.project') as ProjectRouteModel;
  }
}
