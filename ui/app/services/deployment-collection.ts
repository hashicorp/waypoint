import {
  ListDeploymentsRequest,
  ListDeploymentsResponse,
  OperationOrder,
  Ref,
  Deployment,
} from 'waypoint-pb';
import CollectionsService from './collections';

export default class DeploymentCollectionService extends CollectionsService {
  async fetchData(): Promise<Deployment.AsObject[]> {
    var req = new ListDeploymentsRequest();
    req.setWorkspace();

    let appRef = new Ref.Application();
    appRef.setProject(this.applicationObject!.project);
    appRef.setApplication(this.applicationObject!.application);
    req.setApplication(appRef);

    var order = new OperationOrder();
    order.setLimit(3);
    order.setDesc(true);
    req.setOrder(order);

    let resp: ListDeploymentsResponse = await this.api.client.listDeployments(req, this.api.WithMeta());

    // Store on the service
    this.collection = resp.getDeploymentsList().map((d) => d.toObject());

    return this.collection;
  }
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    deploymentResources: DeploymentCollectionService;
  }
}
