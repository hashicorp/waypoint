import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { ListBuildsRequest, ListBuildsResponse } from 'waypoint-pb';
import CurrentApplicationService from 'waypoint/services/current-application';
import CurrentWorkspaceService from 'waypoint/services/current-workspace';
import BuildCollectionService from 'waypoint/services/build-collection';

export default class Builds extends Route {
  @service api!: ApiService;
  @service currentApplication!: CurrentApplicationService;
  @service currentWorkspace!: CurrentWorkspaceService;
  @service buildCollection!: BuildCollectionService;

  beforeModel() {
    this.buildCollection.setup(
      this.currentWorkspace.ref!.toObject(),
      this.currentApplication.ref!.toObject(),
      this
    );
  }
  async model() {
    await this.buildCollection.waitForCollection();

    return this.buildCollection.collection;
  }
}
