import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import ApiService from 'waypoint/services/api';
import { ListBuildsRequest, ListBuildsResponse, OperationOrder, Ref, Build } from 'waypoint-pb';
import CurrentWorkspaceService from 'waypoint/services/current-workspace';
import BuildCollectionService from 'waypoint/services/build-collection';
import { alias } from '@ember/object/computed';

interface AppMetaCardBuildArgs {
  application: Ref.Application.AsObject;
}

export default class AppMetaCardBuilds extends Component<AppMetaCardBuildArgs> {
  @service api!: ApiService;
  @service currentWorkspace!: CurrentWorkspaceService;
  @service buildCollection!: BuildCollectionService;

  constructor(owner: any, args: any) {
    super(owner, args);
    let { application } = this.args;

    this.buildCollection.setup(this.currentWorkspace.ref!.toObject(), application);
  }

  @alias('buildCollection.collection') collection!: Build.AsObject[];

  get firstBuild(): Build.AsObject | undefined {
    if (this.collection) {
      return this.collection.slice(0, 1)[0];
    }
    return;
  }

  get extraBuilds(): Build.AsObject[] | undefined {
    if (this.collection) {
      return this.collection.slice(1, 3);
    }
    return;
  }
}
