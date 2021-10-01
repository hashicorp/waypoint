import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { StatusReport } from 'waypoint-pb';

interface DockerImageBadgeArgs {
  statusReport: StatusReport.AsObject;
}

export default class DockerImageBadge extends Component<DockerImageBadgeArgs> {
  @service api!: ApiService;

  get states(): unknown {
    return this.args.statusReport.resourcesList.map((r) => JSON.parse(r.stateJson ?? '{}'));
  }

  get imageRefs(): ImageRef[] {
    return findImageRefs(this.states);
  }
}

function findImageRefs(obj: unknown, result: ImageRef[] = []): ImageRef[] {
  if (typeof obj !== 'object') {
    return result;
  }

  if (obj === null) {
    return result;
  }

  for (let [key, value] of Object.entries(obj)) {
    if (key.toLowerCase() === 'image' && typeof value === 'string') {
      result.push(new ImageRef(value));
    } else {
      findImageRefs(value, result);
    }
  }

  return result;
}

class ImageRef {
  ref: string;

  constructor(ref: string) {
    this.ref = ref;
  }

  get label(): string {
    return this.split[0];
  }

  get tag(): string {
    return this.split[1];
  }

  private get split(): string[] {
    return this.ref.split(':');
  }
}
