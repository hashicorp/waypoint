import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import ApiService from 'waypoint/services/api';
import { Build } from 'waypoint-pb';
import CurrentWorkspaceService from 'waypoint/services/current-workspace';
import { tracked } from '@glimmer/tracking';

interface AppMetaCardBuildArgs {
  builds: Promise<Build.AsObject[]>;
}

export default class AppMetaCardBuilds extends Component<AppMetaCardBuildArgs> {
  @tracked builds!: Build.AsObject[];
  @tracked loaded!: Boolean;

  constructor(owner: any, args: any) {
    super(owner, args);
    this.load();
  }

  async load() {
    this.builds = await this.args.builds;
    this.loaded = true;
  }

  get firstBuild(): Build.AsObject | undefined {
    return this.builds.slice(0, 1)[0];
  }

  get extraBuilds(): Build.AsObject[] | undefined {
    return this.builds.slice(1, 3);
  }
}
