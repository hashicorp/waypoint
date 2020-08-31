import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { Release } from 'waypoint-pb';

interface LatestReleaseUrlArgs {
  releases: Release.AsObject[];
}

export default class LatestReleaseUrl extends Component<LatestReleaseUrlArgs> {
  @tracked releases!: Release.AsObject[];

  constructor(owner: any, args: any) {
    super(owner, args);
    this.load();
  }

  async load() {
    this.releases = await this.args.releases;
  }

  get firstRelease(): Release.AsObject | undefined {
    if (this.releases) {
      return this.releases[0];
    }
    return;
  }
}
