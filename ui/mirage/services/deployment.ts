import { GetDeploymentRequest, ListDeploymentsRequest, ListDeploymentsResponse } from 'waypoint-pb';
import { Request, Response } from 'ember-cli-mirage';
import { decode } from '../helpers/protobufs';

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function list(schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(ListDeploymentsRequest, requestBody);
  let projectName = requestMsg.getApplication().getProject();
  let appName = requestMsg.getApplication().getApplication();
  let workspaceName = requestMsg.getWorkspace().getWorkspace();
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
export function get(schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(GetDeploymentRequest, requestBody);
  let id = requestMsg.getRef().getId();
  let model = schema.deployments.find(id);
  let protobuf = model?.toProtobuf();

  return this.serialize(protobuf, 'application');
}
