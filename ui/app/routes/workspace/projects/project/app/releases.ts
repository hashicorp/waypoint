/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';
import { Model as AppRouteModel } from '../app';

type Model = {
  app: string;
  releases: AppRouteModel['releases'];
  releaseDeploymentPairs: Record<number, number>;
};

export default class Releases extends Route {
  breadcrumbs(model: Model): Breadcrumb[] {
    if (!model) return [];
    return [
      {
        label: model.app ?? 'unknown',
        route: 'workspace.projects.project.app',
      },
      {
        label: 'Releases',
        route: 'workspace.projects.project.app.releases',
      },
    ];
  }

  async model(): Promise<Model> {
    let app = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let rdPairs = {};
    for (let release of app.releases) {
      let matchingDeployment = app.deployments.find((obj) => obj.id === release?.deploymentId);
      if (matchingDeployment) {
        rdPairs[release.sequence] = matchingDeployment.sequence;
      }
    }

    return {
      app: app.application.application,
      releases: app.releases,
      releaseDeploymentPairs: rdPairs,
    };
  }
}
