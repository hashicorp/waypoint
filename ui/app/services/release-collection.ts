import { ListReleasesRequest, ListReleasesResponse, OperationOrder, Ref, Release } from 'waypoint-pb';
import CollectionsService from './collections';

export default class ReleaseCollectionService extends CollectionsService {
  async fetchData(): Promise<Release.AsObject[]> {
    var req = new ListReleasesRequest();
    req.setWorkspace();

    let appRef = new Ref.Application();
    appRef.setProject(this.applicationObject!.project);
    appRef.setApplication(this.applicationObject!.application);
    req.setApplication(appRef);

    var order = new OperationOrder();
    order.setLimit(3);
    order.setDesc(true);
    req.setOrder(order);

    let resp: ListReleasesResponse = await this.api.client.listReleases(req, this.api.WithMeta());

    // Store on the service
    this.collection = resp.getReleasesList().map((d) => d.toObject());

    return this.collection;
  }
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    releaseResources: ReleaseCollectionService;
  }
}
