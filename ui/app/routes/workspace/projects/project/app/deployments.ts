import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import CurrentApplicationService from 'waypoint/services/current-application';
import CurrentWorkspaceService from 'waypoint/services/current-workspace';
import DeploymentCollectionService from 'waypoint/services/deployment-collection';

export default class Deployments extends Route {
  @service api!: ApiService;
  @service currentApplication!: CurrentApplicationService;
  @service currentWorkspace!: CurrentWorkspaceService;
  @service deploymentCollection!: DeploymentCollectionService;

  beforeModel() {
    this.deploymentCollection.setup(
      this.currentWorkspace.ref!.toObject(),
      this.currentApplication.ref!.toObject(),
      this
    );
  }

  async model() {
    await this.deploymentCollection.waitForCollection();

    return this.deploymentCollection.collection;
  }
}
