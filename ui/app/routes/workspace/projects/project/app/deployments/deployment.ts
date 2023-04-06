/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';
import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import Transition from '@ember/routing/-private/transition';
import ApiService from 'waypoint/services/api';
import { Model as AppRouteModel } from '../../app';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';
import { DeploymentExtended, ReleaseExtended } from 'waypoint/services/api';
import { TimelineModel } from '../../../../../../components/timeline';

type Params = { sequence: string };
export type Model = DeploymentExtended & WithRelease & WithTimeline;

interface WithRelease {
  release?: ReleaseExtended;
}
interface WithTimeline {
  timeline: TimelineModel;
}
export default class DeploymentDetail extends Route {
  @service api!: ApiService;

  breadcrumbs(model: Model): Breadcrumb[] {
    if (!model) return [];
    return [
      {
        label: model.application?.application ?? 'unknown',
        route: 'workspace.projects.project.app',
      },
      {
        label: 'Deployments',
        route: 'workspace.projects.project.app.deployments',
      },
      {
        label: `v${model.sequence}`,
        route: 'workspace.projects.project.app.deployments.deployment',
      },
    ];
  }

  async model(params: Params): Promise<Model> {
    let { builds, deployments, releases } = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let deployment = deployments.find((obj) => obj.sequence == Number(params.sequence));
    let deploymentArtifactId = deployment?.pushedArtifact?.id ?? deployment?.artifactId;
    let build = builds.find((obj) => obj.pushedArtifact?.id === deploymentArtifactId);

    if (!deployment) {
      throw new Error(`Deployment v${params.sequence} not found`);
    }

    let deploymentId = deployment.id;
    let release = releases.find((r) => r.deploymentId === deploymentId);

    let timeline: TimelineModel = {};
    if (build) {
      timeline.build = {
        sequence: build.sequence,
        status: build.status,
      };
    }

    timeline.deployment = {
      sequence: deployment.sequence,
      status: deployment.status,
    };

    if (release) {
      timeline.release = {
        sequence: release.sequence,
        status: release.status,
      };
    }

    return { ...deployment, release, timeline };
  }

  resetController(_: Controller, isExiting: boolean, transition: Transition): void {
    if (isExiting && transition.to.name === 'workspace.projects.project.app.deployment.index') {
      this.transitionTo(
        'workspace.projects.project.app.deployments.deployment',
        (this.modelFor(this.routeName) as Model).sequence
      );
    }
  }
}
