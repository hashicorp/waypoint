import { ListBuildsRequest, ListBuildsResponse, GetBuildRequest } from 'waypoint-pb';
import { Request, Response } from 'miragejs';
import { RouteHandler, Schema } from '../types';
import { decode } from '../helpers/protobufs';

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function list(this: RouteHandler, schema: Schema, { requestBody }: Request): Response {
  let requestMsg = decode(ListBuildsRequest, requestBody);
  let projectName = requestMsg.getApplication()?.getProject();
  let appName = requestMsg.getApplication()?.getApplication();
  let workspaceName = requestMsg.getWorkspace()?.getWorkspace();
  let project = schema.findBy('project', { name: projectName });
  if (!project) {
    throw `Project ${projectName} not found`;
  }
  let application = schema.findBy('application', { name: appName, projectId: project.id });
  let workspace = schema.findBy('workspace', { name: workspaceName });
  let builds = schema.where('build', { applicationId: application?.id, workspaceId: workspace?.id });
  let buildProtobufs = builds.models.map((b) => b.toProtobuf());
  let resp = new ListBuildsResponse();

  buildProtobufs.sort((a, b) => b.getSequence() - a.getSequence());

  resp.setBuildsList(buildProtobufs);

  return this.serialize(resp, 'application');
}

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function get(this: RouteHandler, schema: Schema, { requestBody }: Request): Response {
  let requestMsg = decode(GetBuildRequest, requestBody);
  let id = requestMsg.getRef()?.getId();
  if (!id) {
    throw 'id is required';
  }
  let model = schema.find('build', id);
  let build = model?.toProtobuf();

  return this.serialize(build, 'application');
}
