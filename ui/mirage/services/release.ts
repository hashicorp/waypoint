import { ListReleasesRequest, ListReleasesResponse, GetReleaseRequest } from 'waypoint-pb';
import { Request, Response } from 'ember-cli-mirage';
import { decode } from '../helpers/protobufs';

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function list(schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(ListReleasesRequest, requestBody);
  let projectName = requestMsg.getApplication().getProject();
  let appName = requestMsg.getApplication().getApplication();
  let workspaceName = requestMsg.getWorkspace().getWorkspace();
  let project = schema.projects.findBy({ name: projectName });
  let application = schema.applications.findBy({ name: appName, projectId: project?.id });
  let workspace = schema.workspaces.findBy({ name: workspaceName });
  let releases = schema.releases.where({ applicationId: application?.id, workspaceId: workspace?.id });
  let releaseProtobufs = releases.models.map((d) => d.toProtobuf());
  let resp = new ListReleasesResponse();

  releaseProtobufs.sort((a, b) => b.getSequence() - a.getSequence());

  resp.setReleasesList(releaseProtobufs);

  return this.serialize(resp, 'application');
}

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function get(schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(GetReleaseRequest, requestBody);
  let id = requestMsg.getRef().getId();
  let model = schema.releases.find(id);
  let protobuf = model?.toProtobuf();

  return this.serialize(protobuf, 'application');
}
