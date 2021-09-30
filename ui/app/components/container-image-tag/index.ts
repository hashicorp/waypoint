import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { StatusReport } from 'waypoint-pb';
import { tracked } from '@glimmer/tracking';

interface DockerImageBadgeArgs {
  statusReport: StatusReport.AsObject;
}

export default class DockerImageBadge extends Component<DockerImageBadgeArgs> {
  @service api!: ApiService;
  @tracked image?: string;
  @tracked tag?: string;
  imageFromJson?: string;

  constructor(owner: unknown, args: DockerImageBadgeArgs) {
    super(owner, args);

    this.parseImageAndTag();
  }

  findImageKey(obj: Record<string, unknown>): void {
    if (typeof obj !== 'object') {
      return;
    }

    for (let k in obj) {
      if (k === 'Config') {
        this.imageFromJson = obj[k].Image;
      }

      if (!obj || typeof obj[k] === 'object') {
        this.findImageKey(obj[k] as Record<string, unknown>);
      }
    }
    return;
  }

  parseImageAndTag(): void {
    if (!this.args.statusReport || !this.args.statusReport.resourcesList) {
      return;
    }

    let container = this.args.statusReport.resourcesList.find((r) => r.type === 'container');
    let containerState = JSON.parse(container?.stateJson ?? '{}');
    this.findImageKey(containerState);
    if (this.imageFromJson) {
      [this.image, this.tag] = this.imageFromJson?.split(':');
    }
  }
}
