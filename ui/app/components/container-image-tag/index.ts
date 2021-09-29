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

  constructor(owner: unknown, args: DockerImageBadgeArgs) {
    super(owner, args);

    this.parseImageAndTag();
  }

  parseImageAndTag(): void {
    let container = this.args.statusReport.resourcesList.find((r) => r.type === 'container');
    let containerState = JSON.parse(container?.stateJson ?? '');
    [this.image, this.tag] = containerState['dockerContainerInfo']['Config']['Image'].split(':');
  }
}
