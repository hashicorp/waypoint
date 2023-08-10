/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';
import { Model as AppRouteModel } from '../app';

type Model = {
  app: string;
  builds: AppRouteModel['builds'];
  buildDeploymentPairs: Record<number, number>;
};

export default class Builds extends Route {
  breadcrumbs(model: Model): Breadcrumb[] {
    if (!model) return [];
    return [
      {
        label: model.app ?? 'unknown',
        route: 'workspace.projects.project.app',
      },
      {
        label: 'Builds',
        route: 'workspace.projects.project.app.builds',
      },
    ];
  }

  async model(): Promise<Model> {
    let app = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let bdPairs = {};
    for (let build of app.builds) {
      let matchingDeployment = app.deployments.find((obj) => {
        if (obj.preload && obj.preload.artifact) {
          return obj.preload.artifact.id === build.pushedArtifact?.id;
        }
        return obj.artifactId === build.pushedArtifact?.id;
      });

      if (matchingDeployment) {
        bdPairs[build.sequence] = matchingDeployment.sequence;
      }
    }

    return {
      app: app.application.application,
      builds: app.builds,
      buildDeploymentPairs: bdPairs,
    };
  }
}
