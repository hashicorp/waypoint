import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import ApiService from 'waypoint/services/api';
import { Ref, Release } from 'waypoint-pb';
import { tracked } from '@glimmer/tracking';

interface AppMetaCardReleaseArgs {
  releases: Release.AsObject[];
}

export default class AppMetaCardReleases extends Component<AppMetaCardReleaseArgs> {
  @tracked releases!: Release.AsObject[];
  @tracked loaded!: Boolean;

  constructor(owner: any, args: any) {
    super(owner, args);
    this.load();
  }

  async load() {
    this.releases = await this.args.releases;
    this.loaded = true;
  }

  get firstRelease(): Release.AsObject | undefined {
    return this.releases.slice(0, 1)[0];
    return;
  }

  get extraReleases(): Release.AsObject[] | undefined {
    return this.releases.slice(1, 3);
  }
}
