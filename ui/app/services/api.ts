import Service from '@ember/service';
import { WaypointClient } from 'waypoint-client';
import SessionService from 'waypoint/services/session';
import { inject as service } from '@ember/service';
import {
  Application,
  Build,
  Deployment,
  GetLatestStatusReportRequest,
  Job,
  ListBuildsRequest,
  ListBuildsResponse,
  ListDeploymentsResponse,
  ListPushedArtifactsRequest,
  ListReleasesRequest,
  ListReleasesResponse,
  ListStatusReportsRequest,
  ListStatusReportsResponse,
  OperationOrder,
  Project,
  PushedArtifact,
  Ref,
  Release,
  StatusReport,
  UpsertProjectRequest,
  Variable,
  UI,
} from 'waypoint-pb';
import { Metadata } from 'grpc-web';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
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
  WithMeta(meta?: Metadata): Metadata {
    // In the future we may want additional metadata per-request so this
    // helper merges that per-request metadata supplied at the client request
    // with our authentication metadata
    return { ...this.meta, ...meta };
  }

  get meta(): Metadata {
    if (this.session.authConfigured) {
      return { ...protocolVersions, authorization: this.session.token };
    } else {
      return { ...protocolVersions };
    }
  }

  async listDeployments(
    wsRef: Ref.Workspace,
    appRef: Ref.Application
  ): Promise<UI.DeploymentBundle.AsObject[]> {
    let req = new UI.ListDeploymentsRequest();
    req.setWorkspace(wsRef);
    req.setApplication(appRef);

    let order = new OperationOrder();
    order.setDesc(true);
    req.setOrder(order);

    let resp: UI.ListDeploymentsResponse = await this.client.uI_ListDeployments(req, this.WithMeta());

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
  ): Promise<StatusReport.AsObject | undefined> {
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

  _populateVariableList(
    variablesList: Variable.AsObject[],
    variable?: Variable.AsObject,
    initialVariable?: Variable.AsObject
  ): Variable[] {
    if (variable && initialVariable) {
      let existingVarIndex = variablesList.findIndex((v) => v.name === initialVariable.name);
      if (existingVarIndex !== -1) {
        variablesList.splice(existingVarIndex, 1, variable);
        variablesList = [...variablesList];
      }
    }

    let varProtosList = variablesList.map((v: Variable.AsObject) => {
      let variable = new Variable();
      variable.setName(v.name);
      variable.setServer(new Empty());
      if (v.hcl) {
        variable.setHcl(v.hcl);
      } else {
        variable.setStr(v.str);
      }
      return variable;
    });
    return varProtosList;
  }

  _checkAuthCase(git: Job.Git.AsObject): number {
    if (git.url) {
      if (git.ssh?.privateKeyPem) {
        return 5;
      }
      if (!git?.basic?.username) {
        return 1;
      }
    }
    return 4;
  }

  async upsertProject(
    project: Project.AsObject,
    newAuthCase = -1,
    variable?: Variable.AsObject,
    initialVariable?: Variable.AsObject,
    editedVariableList?: Variable.AsObject[]
  ): Promise<Project.AsObject | undefined> {
    let ref = new Project();
    ref.setName(project.name);

    // Data source settings
    let dataSource = new Job.DataSource();
    let dataSourcePoll = new Project.Poll();
    if (project.dataSourcePoll) {
      dataSourcePoll.setEnabled(project.dataSourcePoll.enabled);
      dataSourcePoll.setInterval(project.dataSourcePoll.interval);
    }

    let git = new Job.Git();

    // Git settings
    if (project?.dataSource?.git) {
      let projGit = project.dataSource.git;

      git.setUrl(projGit.url);
      git.setPath(projGit.path);
      if (!projGit.ref) {
        git.setRef('HEAD');
      } else {
        git.setRef(projGit.ref);
      }

      // get auth case based on existing project settings
      // but if we give a new auth case to this function,
      // that means we're trying to change the auth settings
      let authCase = this._checkAuthCase(projGit);
      if (newAuthCase >= 0) {
        authCase = newAuthCase;
      }

      // Git authentication settings
      if (authCase === 4) {
        let gitBasic = new Job.Git.Basic();
        gitBasic.setUsername(projGit.basic?.username ?? '');
        gitBasic.setPassword(projGit.basic?.password ?? '');
        git.setBasic(gitBasic);
        git.clearSsh();
      }

      // SSH authentication settings
      if (authCase === 5) {
        let gitSSH = new Job.Git.SSH();
        gitSSH.setPrivateKeyPem(projGit.ssh?.privateKeyPem ?? '');
        gitSSH.setUser(projGit.ssh?.user ?? '');
        gitSSH.setPassword(projGit.ssh?.password ?? '');
        git.setSsh(gitSSH);
        git.clearBasic();
      }

      // Basic authentication settings
      if (authCase === 0) {
        git.clearBasic();
        git.clearSsh();
      }
    } else {
      // if we set up a project without connecting it to a git repo
      // but we want to set input variables, a git URL is required
      // for updating a project's settings. this silences that error
      // while not adding settings the user did not specify
      git.setUrl('\n');
    }

    dataSource.setGit(git);
    ref.setDataSource(dataSource);
    ref.setDataSourcePoll(dataSourcePoll);

    if (project.waypointHcl) {
      // Hardcode hcl for now
      ref.setWaypointHclFormat(0); // check project-repository-settings.ts for FORMAT obj
      ref.setWaypointHcl(project.waypointHcl);
    }

    // Application list settings
    let appList = project.applicationsList.map(applicationFromObject);
    ref.setApplicationsList(appList);

    // Input variable settings
    let startingList = project.variablesList;
    if (editedVariableList) {
      startingList = editedVariableList;
    }
    let varsList = this._populateVariableList(startingList, variable, initialVariable);
    ref.setVariablesList(varsList);

    // Build and trigger request
    let req = new UpsertProjectRequest();
    req.setProject(ref);

    let resp = await this.client.upsertProject(req, this.WithMeta());
    let respProject = resp.toObject().project;
    return respProject;
  }
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    api: ApiService;
  }
}

function applicationFromObject(object: Application.AsObject): Application {
  let result = new Application();

  result.setName(object.name);
  result.setFileChangeSignal(object.fileChangeSignal);

  if (object.project) {
    let ref = new Ref.Project();
    ref.setProject(object.project.project);
    result.setProject(ref);
  }

  return result;
}
