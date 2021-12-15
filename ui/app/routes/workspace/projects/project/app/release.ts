import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Model as AppRouteModel } from '../app';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';
import { ReleaseExtended } from 'waypoint/services/api';
import { TimelineModel } from '../../../../../components/timeline';

type Params = { sequence: string };

interface WithTimeline {
  timeline: TimelineModel;
}

export type Model = ReleaseExtended & WithTimeline;

export default class ReleaseDetail extends Route {
  @service api!: ApiService;

  breadcrumbs(model: Model): Breadcrumb[] {
    if (!model) return [];

    return [
      {
        label: model.application?.application ?? 'unknown',
        icon: 'git-repo',
        route: 'workspace.projects.project.app',
      },
      {
        label: 'Releases',
        icon: 'globe',
        route: 'workspace.projects.project.app.releases',
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

    return { ...release, timeline };
  }
}
