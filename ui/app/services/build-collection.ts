import { ListBuildsRequest, ListBuildsResponse, OperationOrder, Ref, Build } from 'waypoint-pb';
import CollectionsService from './collections';

export default class BuildCollectionService extends CollectionsService {
  async fetchData(): Promise<Build.AsObject[]> {
    var req = new ListBuildsRequest();
    req.setWorkspace();

    let appRef = new Ref.Application();
    appRef.setProject(this.applicationObject!.project);
    appRef.setApplication(this.applicationObject!.application);
    req.setApplication(appRef);

    var order = new OperationOrder();
    order.setLimit(3);
    order.setDesc(true);
    // todo(pearkes): set order
    // req.setOrder(order);

    let resp: ListBuildsResponse = await this.api.client.listBuilds(req, this.api.WithMeta());

    // Store on the service
    this.collection = resp.getBuildsList().map((d) => d.toObject());

    return this.collection;
  }
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    buildResources: BuildCollectionService;
  }
}
