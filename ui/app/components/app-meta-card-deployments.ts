import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import ApiService from 'waypoint/services/api';
import {
  ListDeploymentsRequest,
  ListDeploymentsResponse,
  OperationOrder,
  Ref,
  Deployment,
} from 'waypoint-pb';
import CurrentWorkspaceService from 'waypoint/services/current-workspace';
import { task } from 'ember-concurrency-decorators';
import { taskFor } from 'ember-concurrency-ts';
import { tracked } from '@glimmer/tracking';

interface AppMetaCardDeploymentArgs {
  application: Ref.Application.AsObject;
}

export default class AppMetaCardDeployments extends Component<AppMetaCardDeploymentArgs> {
  @service api!: ApiService;
  @service currentWorkspace!: CurrentWorkspaceService;

  @task async fetchData(): Promise<Deployment.AsObject[]> {
    var req = new ListDeploymentsRequest();
    req.setWorkspace(this.currentWorkspace.ref);

    let appRef = new Ref.Application();
    appRef.setProject(this.applicationObject!.project);
    appRef.setApplication(this.applicationObject!.application);
    req.setApplication(appRef);

    var order = new OperationOrder();
    order.setLimit(3);
    order.setDesc(true);
    req.setOrder(order);

    let resp: ListDeploymentsResponse = await this.api.client.listDeployments(req, this.api.WithMeta());
    return resp.getDeploymentsList().map((d) => d.toObject());
  }

  performTask() {
    taskFor(this.fetchData)
      .perform()
      .then((deployments) => {
        if (deployments.length > 0) {
          this.firstDeployment = deployments.splice(0, 1)[0];
          this.deployments = deployments;
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
  @tracked deployments?: Deployment.AsObject[];

  @tracked firstDeployment?: Deployment.AsObject;
}
