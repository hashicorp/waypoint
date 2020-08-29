import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import ApiService from 'waypoint/services/api';
import CurrentWorkspaceService from 'waypoint/services/current-workspace';
import { alias } from '@ember/object/computed';
import { Ref, Release } from 'waypoint-pb';
import ReleaseCollectionService from 'waypoint/services/release-collection';

interface AppMetaCardReleaseArgs {
  application: Ref.Application.AsObject;
}

export default class AppMetaCardReleases extends Component<AppMetaCardReleaseArgs> {
  @service api!: ApiService;
  @service currentWorkspace!: CurrentWorkspaceService;
  @service releaseCollection!: ReleaseCollectionService;

  constructor(owner: any, args: any) {
    super(owner, args);
    let { application } = this.args;

    this.releaseCollection.setup(this.currentWorkspace.ref!.toObject(), application);
  }

  @alias('releaseCollection.collection') collection!: Release.AsObject[];

  get firstRelease(): Release.AsObject | undefined {
    if (this.collection) {
      return this.collection.slice(0, 1)[0];
    }
    return;
  }

  get extraReleases(): Release.AsObject[] | undefined {
    if (this.collection) {
      return this.collection.slice(1, 3);
    }
    return;
  }
}
