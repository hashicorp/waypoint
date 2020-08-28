import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import ApiService from 'waypoint/services/api';
import { ListReleasesRequest, ListReleasesResponse, OperationOrder, Ref, Release } from 'waypoint-pb';
import CurrentWorkspaceService from 'waypoint/services/current-workspace';
import { task } from 'ember-concurrency-decorators';
import { taskFor } from 'ember-concurrency-ts';
import { tracked } from '@glimmer/tracking';

interface AppMetaCardReleaseArgs {
  application: Ref.Application.AsObject;
}

export default class AppMetaCardReleases extends Component<AppMetaCardReleaseArgs> {
  @service api!: ApiService;
  @service currentWorkspace!: CurrentWorkspaceService;

  @task async fetchData(): Promise<Release.AsObject[]> {
    var req = new ListReleasesRequest();
    req.setWorkspace(this.currentWorkspace.ref);

    let appRef = new Ref.Application();
    appRef.setProject(this.applicationObject!.project);
    appRef.setApplication(this.applicationObject!.application);
    req.setApplication(appRef);

    var order = new OperationOrder();
    order.setLimit(3);
    order.setDesc(true);
    req.setOrder(order);

    let resp: ListReleasesResponse = await this.api.client.listReleases(req, this.api.WithMeta());
    return resp.getReleasesList().map((d) => d.toObject());
  }

  performTask() {
    ``;
    taskFor(this.fetchData)
      .perform()
      .then((releases) => {
        if (releases.length > 0) {
          this.firstRelease = releases.splice(0, 1)[0];
          this.releases = releases;
        }
      });
  }

  constructor(owner: any, args: any) {
    super(owner, args);

    let { application } = this.args;
    this.applicationObject = application;

    this.performTask();
  }

  @tracked applicationObject?: Ref.Application.AsObject;
  @tracked releases?: Release.AsObject[];

  @tracked firstRelease?: Release.AsObject;
}
