/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Request, Response } from 'miragejs';
import { RouteHandler } from '../types';
import { ListPushedArtifactsRequest, ListPushedArtifactsResponse } from 'waypoint-pb';
import { decode } from '../helpers/protobufs';

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function list(this: RouteHandler, schema: any, request: Request): Response {
  let requestMsg = decode(ListPushedArtifactsRequest, request.requestBody);
  let projectName = requestMsg.getApplication()?.getProject();
  let appName = requestMsg.getApplication()?.getApplication();
  let workspaceName = requestMsg.getWorkspace()?.getWorkspace();
  let project = schema.projects.findBy({ name: projectName });
  let application = schema.applications.findBy({ name: appName, projectId: project.id });
  let workspace = schema.workspaces.findBy({ name: workspaceName });
  let pushedArtifacts = schema.pushedArtifacts.where({
    applicationId: application?.id,
    workspaceId: workspace?.id,
  });
  let pushedArtifactProtobufs = pushedArtifacts.models.map((b) => b.toProtobuf());
  let resp = new ListPushedArtifactsResponse();

  pushedArtifactProtobufs.sort((a, b) => b.getSequence() - a.getSequence());

  resp.setArtifactsList(pushedArtifactProtobufs);

  return this.serialize(resp, 'application');
}
