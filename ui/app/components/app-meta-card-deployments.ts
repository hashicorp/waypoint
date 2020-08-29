import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import ApiService from 'waypoint/services/api';
import { Ref, Deployment } from 'waypoint-pb';
import CurrentWorkspaceService from 'waypoint/services/current-workspace';
import { alias } from '@ember/object/computed';
import DeploymentCollectionService from 'waypoint/services/deployment-collection';

interface AppMetaCardDeploymentArgs {
  application: Ref.Application.AsObject;
}

export default class AppMetaCardDeployments extends Component<AppMetaCardDeploymentArgs> {
  @service api!: ApiService;
  @service currentWorkspace!: CurrentWorkspaceService;
  @service deploymentCollection!: DeploymentCollectionService;

  constructor(owner: any, args: any) {
    super(owner, args);
    let { application } = this.args;

    this.deploymentCollection.setup(this.currentWorkspace.ref!.toObject(), application);
  }

  @alias('deploymentCollection.collection') collection!: Deployment.AsObject[];

  get firstDeployment(): Deployment.AsObject | undefined {
    if (this.collection) {
      return this.collection.slice(0, 1)[0];
    }
    return;
  }

  get extraDeployments(): Deployment.AsObject[] | undefined {
    if (this.collection) {
      return this.collection.slice(1, 3);
    }
    return;
  }
}
