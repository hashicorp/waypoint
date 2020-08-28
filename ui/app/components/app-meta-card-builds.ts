import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import ApiService from 'waypoint/services/api';
import { ListBuildsRequest, ListBuildsResponse, OperationOrder, Ref, Build } from 'waypoint-pb';
import CurrentWorkspaceService from 'waypoint/services/current-workspace';
import { task } from 'ember-concurrency-decorators';
import { taskFor } from 'ember-concurrency-ts';
import { tracked } from '@glimmer/tracking';

interface AppMetaCardBuildArgs {
  application: Ref.Application.AsObject;
}

export default class AppMetaCardBuilds extends Component<AppMetaCardBuildArgs> {
  @service api!: ApiService;
  @service currentWorkspace!: CurrentWorkspaceService;

  @task async fetchData(): Promise<Build.AsObject[]> {
    var req = new ListBuildsRequest();
    req.setWorkspace(this.currentWorkspace.ref);

    let appRef = new Ref.Application();
    appRef.setProject(this.applicationObject!.project);
    appRef.setApplication(this.applicationObject!.application);
    req.setApplication(appRef);

    var order = new OperationOrder();
    order.setLimit(3);
    order.setDesc(true);
    // todo(pearkes): builds need order set
    // req.setOrder(order);

    let resp: ListBuildsResponse = await this.api.client.listBuilds(req, this.api.WithMeta());
    return resp.getBuildsList().map((d) => d.toObject());
  }

  performTask() {
    taskFor(this.fetchData)
      .perform()
      .then((builds) => {
        if (builds.length > 0) {
          this.firstBuild = builds.splice(0, 1)[0];
          this.builds = builds;
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
  @tracked builds?: Build.AsObject[];

  @tracked firstBuild?: Build.AsObject;
}
