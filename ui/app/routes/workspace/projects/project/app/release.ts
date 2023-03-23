/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Model as AppRouteModel } from '../app';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';
import { ReleaseExtended } from 'waypoint/services/api';
import { TimelineModel } from '../../../../../components/timeline';
import { Operation } from 'waypoint-pb';

type Params = { sequence: string };

interface WithTimeline {
  timeline: TimelineModel;
}

interface WithLatest {
  latestRelease?: ReleaseExtended;
}

export type Model = ReleaseExtended & WithTimeline & WithLatest;

export default class ReleaseDetail extends Route {
  @service api!: ApiService;

  breadcrumbs(model: Model): Breadcrumb[] {
    if (!model) return [];

    return [
      {
        label: model.application?.application ?? 'unknown',
        route: 'workspace.projects.project.app',
      },
      {
        label: 'Releases',
        route: 'workspace.projects.project.app.releases',
      },
      {
        label: `v${model.sequence}`,
        route: 'workspace.projects.project.app.release',
      },
    ];
  }

  async model(params: Params): Promise<Model> {
    let { builds, deployments, releases } = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let release = releases.find((obj) => obj.sequence === Number(params.sequence));

    if (!release) {
      throw new Error(`Release v${params.sequence} not found`);
    }

    let deployment = deployments.find((obj) => obj.id === release?.deploymentId);
    let deploymentArtifactId = deployment?.pushedArtifact?.id ?? deployment?.artifactId;
    let build = builds.find((obj) => obj.pushedArtifact?.id === deploymentArtifactId);
    let latestRelease = releases.find((r) => r.state === Operation.PhysicalState.CREATED);

    let timeline: TimelineModel = {};
    if (build) {
      let buildObj = {
        sequence: build.sequence,
        status: build.status,
      };
      timeline.build = buildObj;
    }

    if (deployment) {
      timeline.deployment = {
        sequence: deployment.sequence,
        status: deployment.status,
      };
    }

    let releaseObj = {
      sequence: release.sequence,
      status: release.status,
    };
    timeline.release = releaseObj;

    return { ...release, timeline, latestRelease };
  }
}
