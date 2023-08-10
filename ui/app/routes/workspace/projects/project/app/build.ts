/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Ref, GetBuildRequest, Build, PushedArtifact } from 'waypoint-pb';
import { Model as AppRouteModel } from '../app';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';
import { TimelineModel } from 'waypoint/components/timeline';

type Params = { sequence: string };
type Model = Build.AsObject & WithPushedArtifact;

interface WithPushedArtifact {
  pushedArtifact?: PushedArtifact.AsObject;
}

interface WithTimeline {
  timeline: TimelineModel;
}

type BuildWithArtifact = Build.AsObject & WithPushedArtifact;
type BuildWithArtifactAndTimeline = BuildWithArtifact & WithTimeline;

export default class BuildDetail extends Route {
  @service api!: ApiService;

  breadcrumbs(model: Model): Breadcrumb[] {
    if (!model) return [];
    return [
      {
        label: model.application?.application ?? 'unknown',
        route: 'workspace.projects.project.app',
      },
      {
        label: 'Builds',
        route: 'workspace.projects.project.app.builds',
      },
      {
        label: `v${model.sequence}`,
        route: 'workspace.projects.project.app.build',
      },
    ];
  }

  async model(params: Params): Promise<Model> {
    let { builds, deployments, releases } = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let buildFromAppRoute = builds.find((obj) => obj.sequence === Number(params.sequence));
    let deploymentFromAppRoute = deployments.find((obj) => {
      if (obj.preload && obj.preload.artifact) {
        return obj.preload.artifact.id === buildFromAppRoute?.pushedArtifact?.id;
      }
      return obj.artifactId === buildFromAppRoute?.pushedArtifact?.id;
    });
    let releaseFromAppRoute = releases.find((obj) => obj.deploymentId === deploymentFromAppRoute?.id);

    if (!buildFromAppRoute) {
      throw new Error(`Build v${params.sequence} not found`);
    }

    let ref = new Ref.Operation();
    ref.setId(buildFromAppRoute.id);
    let req = new GetBuildRequest();
    req.setRef(ref);

    let build = await this.api.client.getBuild(req, this.api.WithMeta());
    let buildWithArtifact: BuildWithArtifact = build.toObject();

    buildWithArtifact.pushedArtifact = buildFromAppRoute.pushedArtifact;

    let timeline: TimelineModel = {};
    timeline.build = {
      sequence: buildFromAppRoute.sequence,
      status: buildFromAppRoute.status,
    };

    if (deploymentFromAppRoute) {
      timeline.deployment = {
        sequence: deploymentFromAppRoute.sequence,
        status: deploymentFromAppRoute.status,
      };
    }

    if (releaseFromAppRoute) {
      timeline.release = {
        sequence: releaseFromAppRoute.sequence,
        status: releaseFromAppRoute.status,
      };
    }

    let result: BuildWithArtifactAndTimeline = { ...buildWithArtifact, timeline };
    return result;
  }
}
