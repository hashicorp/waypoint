import Service from '@ember/service';
import { WaypointClient } from 'waypoint-client';
import SessionService from 'waypoint/services/session';
import { inject as service } from '@ember/service';
import { assign } from '@ember/polyfills';
import {
  ListDeploymentsRequest,
  Ref,
  Deployment,
  OperationOrder,
  ListDeploymentsResponse,
  ListBuildsRequest,
  Build,
  ListBuildsResponse,
  Release,
  ListReleasesRequest,
  ListReleasesResponse,
} from 'waypoint-pb';

export default class ApiService extends Service {
  @service session!: SessionService;
  client = new WaypointClient('/grpc', null, null);

  // Merges metadata with required metadata for the request
  WithMeta(meta?: any) {
    // In the future we may want additional metadata per-request so this
    // helper merges that per-request metadata supplied at the client request
    // with our authentication metadata
    return assign(this.meta, meta!).valueOf();
  }

  get meta() {
    if (this.session.authConfigured) {
      return { authorization: this.session.token };
    } else {
      return {};
    }
  }

  async listDeployments(wsRef: Ref.Workspace, appRef: Ref.Application): Promise<Deployment.AsObject[]> {
    var req = new ListDeploymentsRequest();
    req.setWorkspace(wsRef);
    req.setApplication(appRef);

    var order = new OperationOrder();
    order.setDesc(true);
    req.setOrder(order);

    let resp: ListDeploymentsResponse = await this.client.listDeployments(req, this.WithMeta());

    return resp.getDeploymentsList().map((d) => d.toObject());
  }

  async listBuilds(wsRef: Ref.Workspace, appRef: Ref.Application): Promise<Build.AsObject[]> {
    var req = new ListBuildsRequest();
    req.setWorkspace(wsRef);
    req.setApplication(appRef);

    var order = new OperationOrder();
    order.setLimit(3);
    order.setDesc(true);
    // todo(pearkes): set order
    // req.setOrder(order);

    let resp: ListBuildsResponse = await this.client.listBuilds(req, this.WithMeta());

    return resp.getBuildsList().map((d) => d.toObject());
  }

  async listReleases(wsRef: Ref.Workspace, appRef: Ref.Application): Promise<Release.AsObject[]> {
    var req = new ListReleasesRequest();
    req.setWorkspace(wsRef);
    req.setApplication(appRef);

    var order = new OperationOrder();
    order.setLimit(3);
    order.setDesc(true);
    req.setOrder(order);

    let resp: ListReleasesResponse = await this.client.listReleases(req, this.WithMeta());

    return resp.getReleasesList().map((d) => d.toObject());
  }
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    api: ApiService;
  }
}
