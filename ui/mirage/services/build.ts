/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { ListBuildsRequest, ListBuildsResponse, GetBuildRequest } from 'waypoint-pb';
import { Request, Response } from 'miragejs';
import { RouteHandler } from '../types';
import { decode } from '../helpers/protobufs';

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function list(this: RouteHandler, schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(ListBuildsRequest, requestBody);
  let projectName = requestMsg.getApplication()?.getProject();
  let appName = requestMsg.getApplication()?.getApplication();
  let workspaceName = requestMsg.getWorkspace()?.getWorkspace();
  let project = schema.projects.findBy({ name: projectName });
  let application = schema.applications.findBy({ name: appName, projectId: project.id });
  let workspace = schema.workspaces.findBy({ name: workspaceName });
  let builds = schema.builds.where({ applicationId: application?.id, workspaceId: workspace?.id });
  let buildProtobufs = builds.models.map((b) => b.toProtobuf());
  let resp = new ListBuildsResponse();

  buildProtobufs.sort((a, b) => b.getSequence() - a.getSequence());

  resp.setBuildsList(buildProtobufs);

  return this.serialize(resp, 'application');
}

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function get(this: RouteHandler, schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(GetBuildRequest, requestBody);
  let id = requestMsg.getRef()?.getId();
  let model = schema.builds.find(id);
  let build = model?.toProtobuf();

  return this.serialize(build, 'application');
}
