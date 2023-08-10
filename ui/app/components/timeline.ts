/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import RouterService from '@ember/routing/router-service';
import ApiService from 'waypoint/services/api';
import { Status } from 'waypoint-pb';

const _TYPE_TRANSLATIONS = {
  build: 'page.artifact.timeline.build',
  deployment: 'page.artifact.timeline.deployment',
  release: 'page.artifact.timeline.release',
};

const _TYPE_ROUTES = {
  build: 'workspace.projects.project.app.build',
  deployment: 'workspace.projects.project.app.deployments.deployment',
  release: 'workspace.projects.project.app.release',
};
interface ArtifactModel {
  sequence: number;
  type: string;
  route: string;
  status?: Status.AsObject;
  isCurrentRoute: boolean;
}

export interface TimelineModel {
  build?: TimelineArtifact;
  deployment?: TimelineArtifact;
  release?: TimelineArtifact;
}
interface TimelineArtifact {
  sequence: number;
  status: Status.AsObject | undefined;
}
interface TimelineArgs {
  model: TimelineModel;
}

export default class Timeline extends Component<TimelineArgs> {
  @service api!: ApiService;
  @service router!: RouterService;

  areWeHere(currentArtifactKey: string): boolean {
    let routeForKey = _TYPE_ROUTES[currentArtifactKey];
    return this.router.currentRouteName.includes(routeForKey);
  }

  get artifacts(): ArtifactModel[] {
    let artifactsList: ArtifactModel[] = [];
    for (let key in this.args.model) {
      artifactsList.push({
        sequence: this.args.model[key].sequence,
        type: _TYPE_TRANSLATIONS[key],
        route: _TYPE_ROUTES[key],
        status: this.args.model[key].status,
        isCurrentRoute: this.areWeHere(key),
      } as ArtifactModel);
    }
    artifactsList.sort((a, b) => (a.type > b.type ? 1 : -1));
    return artifactsList;
  }
}
