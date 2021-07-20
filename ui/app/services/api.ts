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
  StatusReport,
  ListStatusReportsRequest,
  ListStatusReportsResponse,
  GetLatestStatusReportRequest,
  ListPushedArtifactsRequest,
  PushedArtifact,
} from 'waypoint-pb';
import config from 'waypoint/config/environment';

const protocolVersions = {
  // These map to upstream protocol versions
  'client-api-protocol': '1,1',
  'client-entrypoint-protocol': '1,1',
  // This is defined by the UI and can be
  // later used to identify different versions of the UI
  // todo: policy for when we change this..
  'client-version': 'ui-0.0.1',
};

export default class ApiService extends Service {
  @service session!: SessionService;
  // If the the apiAddress is not set, this will use the /grpc prefix on the
  // same host as the UI is being served from
  client = new WaypointClient(`${config.apiAddress}/grpc`, null, null);

  // Merges metadata with required metadata for the request
  WithMeta(meta?: any) {
    // In the future we may want additional metadata per-request so this
    // helper merges that per-request metadata supplied at the client request
    // with our authentication metadata
    return assign(this.meta, meta!).valueOf();
  }

  get meta() {
    if (this.session.authConfigured) {
      return { ...protocolVersions, authorization: this.session.token };
    } else {
      return { ...protocolVersions };
    }
  }

  async listDeployments(wsRef: Ref.Workspace, appRef: Ref.Application): Promise<Deployment.AsObject[]> {
    let req = new ListDeploymentsRequest();
    req.setWorkspace(wsRef);
    req.setApplication(appRef);

    let order = new OperationOrder();
    order.setDesc(true);
    req.setOrder(order);

    let resp: ListDeploymentsResponse = await this.client.listDeployments(req, this.WithMeta());

    return resp.getDeploymentsList().map((d) => d.toObject());
  }

  async listBuilds(wsRef: Ref.Workspace, appRef: Ref.Application): Promise<Build.AsObject[]> {
    let req = new ListBuildsRequest();
    req.setWorkspace(wsRef);
    req.setApplication(appRef);

    let order = new OperationOrder();
    order.setLimit(3);
    order.setDesc(true);
    // todo(pearkes): set order
    // req.setOrder(order);

    let resp: ListBuildsResponse = await this.client.listBuilds(req, this.WithMeta());

    return resp.getBuildsList().map((d) => d.toObject());
  }

  async listPushedArtifacts(
    wsRef: Ref.Workspace,
    appRef: Ref.Application
  ): Promise<PushedArtifact.AsObject[]> {
    let request = new ListPushedArtifactsRequest();

    request.setApplication(appRef);
    request.setWorkspace(wsRef);

    // TODO(jgwhite): request.setIncludeBuild
    // TODO(jgwhite): request.setOrder
    // TODO(jgwhite): request.setStatusList

    let response = await this.client.listPushedArtifacts(request, this.WithMeta());
    let result = response.getArtifactsList().map((pa) => pa.toObject());

    return result;
  }

  async listReleases(wsRef: Ref.Workspace, appRef: Ref.Application): Promise<Release.AsObject[]> {
    let req = new ListReleasesRequest();
    req.setWorkspace(wsRef);
    req.setApplication(appRef);

    let order = new OperationOrder();
    order.setLimit(3);
    order.setDesc(true);
    req.setOrder(order);

    let resp: ListReleasesResponse = await this.client.listReleases(req, this.WithMeta());

    return resp.getReleasesList().map((d) => d.toObject());
  }

  async listStatusReports(wsRef: Ref.Workspace, appRef: Ref.Application): Promise<StatusReport.AsObject[]> {
    let req = new ListStatusReportsRequest();
    req.setWorkspace(wsRef);
    req.setApplication(appRef);

    let order = new OperationOrder();
    order.setDesc(true);
    req.setOrder(order);

    let resp: ListStatusReportsResponse = await this.client.listStatusReports(req, this.WithMeta());

    return resp.getStatusReportsList().map((d) => d.toObject());
  }

  async getLatestStatusReport(
    _wsRef: Ref.Workspace,
    appRef: Ref.Application
  ): Promise<StatusReport | undefined> {
    let req = new GetLatestStatusReportRequest();
    req.setApplication(appRef);
    // We have to try/catch to avoid failing the hash request because the api errors if no statusReport is available
    try {
      let resp: StatusReport = await this.client.getLatestStatusReport(req, this.WithMeta());
      return resp.toObject();
    } catch {
      return;
    }
  }
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    api: ApiService;
  }
}
