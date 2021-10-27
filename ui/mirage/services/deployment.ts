import { GetDeploymentRequest, ListDeploymentsRequest, ListDeploymentsResponse, UI } from 'waypoint-pb';
import { RouteHandler, Request, Response } from 'ember-cli-mirage';
import { decode } from '../helpers/protobufs';

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function list(this: RouteHandler, schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(ListDeploymentsRequest, requestBody);
  let projectName = requestMsg.getApplication()?.getProject();
  let appName = requestMsg.getApplication()?.getApplication();
  let workspaceName = requestMsg.getWorkspace()?.getWorkspace();
  let project = schema.projects.findBy({ name: projectName });
  let application = schema.applications.findBy({ name: appName, projectId: project?.id });
  let workspace = schema.workspaces.findBy({ name: workspaceName });
  let deployments = schema.deployments.where({ applicationId: application?.id, workspaceId: workspace?.id });
  let deploymentProtobufs = deployments.models.map((d) => d.toProtobuf());
  let resp = new ListDeploymentsResponse();

  deploymentProtobufs.sort((a, b) => b.getSequence() - a.getSequence());

  resp.setDeploymentsList(deploymentProtobufs);

  return this.serialize(resp, 'application');
}

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function ui_list(this: RouteHandler, schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(UI.ListDeploymentsRequest, requestBody);
  let projectName = requestMsg.getApplication()?.getProject();
  let appName = requestMsg.getApplication()?.getApplication();
  let workspaceName = requestMsg.getWorkspace()?.getWorkspace();
  let project = schema.projects.findBy({ name: projectName });
  let application = schema.applications.findBy({ name: appName, projectId: project?.id });
  let workspace = schema.workspaces.findBy({ name: workspaceName });
  let deployments = schema.deployments
    .where({ applicationId: application?.id, workspaceId: workspace?.id })
    .models.sort((a, b) => b.sequence - a.sequence);
  let bundles: UI.DeploymentBundle[] = deployments.map((deployment) => {
    let bundle = new UI.DeploymentBundle();

    bundle.setDeployment(deployment.toProtobuf());
    bundle.setArtifact(deployment.build?.pushedArtifact?.toProtobuf());
    bundle.setBuild(deployment.build?.toProtobuf());
    bundle.setDeployUrl(deployment.deployUrl);
    bundle.setLatestStatusReport(deployment.statusReport?.toProtobuf());
    // TODO(jgwhite): bundle.setJobDataSourceRef

    return bundle;
  });
  let resp = new UI.ListDeploymentsResponse();

  resp.setDeploymentsList(bundles);

  return this.serialize(resp, 'application');
}

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function get(this: RouteHandler, schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(GetDeploymentRequest, requestBody);
  let id = requestMsg.getRef()?.getId();
  let model = schema.deployments.find(id);
  let protobuf = model?.toProtobuf();

  return this.serialize(protobuf, 'application');
}
