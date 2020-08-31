import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import ApiService from 'waypoint/services/api';
import { Deployment } from 'waypoint-pb';
import { tracked } from '@glimmer/tracking';

interface AppMetaCardDeploymentArgs {
  deployments: Promise<Deployment.AsObject[]>;
}

export default class AppMetaCardDeployments extends Component<AppMetaCardDeploymentArgs> {
  @service api!: ApiService;

  @tracked deployments!: Deployment.AsObject[];
  @tracked loaded!: Boolean;

  constructor(owner: any, args: any) {
    super(owner, args);
    this.load();
  }

  async load() {
    this.deployments = await this.args.deployments;
    this.loaded = true;
  }

  get firstDeployment(): Deployment.AsObject | undefined {
    return this.deployments.slice(0, 1)[0];
  }

  get extraDeployments(): Deployment.AsObject[] | undefined {
    return this.deployments.slice(1, 3);
  }
}
