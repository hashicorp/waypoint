/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { DeploymentExtended, ReleaseExtended } from 'waypoint/services/api';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';
import { Model as AppRouteModel } from '../app';
import Route from '@ember/routing/route';
import { StatusReport } from 'waypoint-pb';

type Model = {
  resources: ResourceMap[];
  application: string;
};

interface ResourceMap {
  resource: StatusReport.Resource.AsObject;
  type: string;
  source: DeploymentExtended | ReleaseExtended;
}
export default class Resources extends Route {
  breadcrumbs(model: Model): Breadcrumb[] {
    if (!model) return [];
    return [
      {
        label: model.application ?? 'unknown',
        route: 'workspace.projects.project.app',
      },
      {
        label: 'Resources',
        route: 'workspace.projects.project.app.resources',
      },
    ];
  }

  async model(): Promise<Model> {
    let app = this.modelFor('workspace.projects.project.app') as AppRouteModel;

    let deployments = app.deployments;
    let releases = app.releases;

    let resources: ResourceMap[] = [];

    deployments.forEach((dep) => {
      dep.statusReport?.resourcesList.forEach((resource) => {
        resources.push({
          resource,
          type: 'deployment',
          source: dep,
        } as ResourceMap);
      });
    });
    releases.forEach((rel) => {
      rel.statusReport?.resourcesList.forEach((resource) => {
        resources.push({
          resource,
          type: 'release',
          source: rel,
        } as ResourceMap);
      });
    });

    return { application: app.application.application, resources };
  }
}
