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
  deployment: 'workspace.projects.project.app.deployment.deployment-seq',
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
    let entry = Object.entries(_TYPE_ROUTES).find(([_, value]) =>
      this.router.currentRouteName.includes(value)
    );

    if (entry) {
      return entry[0] === currentArtifactKey;
    }
    return false;
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
    return artifactsList;
  }
}
