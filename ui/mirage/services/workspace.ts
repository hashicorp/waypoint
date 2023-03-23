/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Request, Response } from 'miragejs';
import { ListWorkspacesRequest, ListWorkspacesResponse, Ref } from 'waypoint-pb';
import { decode } from '../helpers/protobufs';

// eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/explicit-module-boundary-types
export function list(schema: any, request: Request): Response {
  let requestMsg = decode(ListWorkspacesRequest, request.requestBody);
  let response = new ListWorkspacesResponse();
  let workspaces = [];

  switch (requestMsg.getScopeCase()) {
    case ListWorkspacesRequest.ScopeCase.GLOBAL:
    case ListWorkspacesRequest.ScopeCase.SCOPE_NOT_SET:
      workspaces = schema.all('workspace').models;
      break;
    case ListWorkspacesRequest.ScopeCase.PROJECT:
      workspaces = workspacesForProject(schema, requestMsg.getProject());
      break;
    case ListWorkspacesRequest.ScopeCase.APPLICATION:
      workspaces = workspacesForApplication(schema, requestMsg.getApplication());
      break;
  }

  workspaces.sort((a, b) => a.name.localeCompare(b.name));

  response.setWorkspacesList(workspaces.map((w) => w.toProtobuf()));

  return this.serialize(response, 'application');
}

function workspacesForProject(schema, ref: Ref.Project) {
  let name = ref.getProject();
  let project = schema.findBy('project', { name });
  if (!project) {
    throw `Project ${name} not found`;
  }
  let ids = new Set();

  for (let app of project.applications.models) {
    let operations = [...app.builds.models, ...app.deployments.models, ...app.releases.models];
    for (let op of operations) {
      ids.add(op.workspaceId);
    }
  }

  let result = schema.find('workspace', [...ids]).models;

  return result;
}

function workspacesForApplication(schema, ref: Ref.Application) {
  let projectName = ref.getProject();
  let appName = ref.getApplication();
  let project = schema.findBy('project', { name: projectName });
  if (!project) {
    throw `Project ${projectName} not found`;
  }
  let app = schema.findBy('application', { name: appName, projectId: project.id });
  if (!app) {
    throw `Application ${projectName}/${appName} not found`;
  }
  let ids = new Set();

  let operations = [...app.builds.models, ...app.deployments.models, ...app.releases.models];

  for (let op of operations) {
    ids.add(op.workspaceId);
  }

  let result = schema.find('workspace', [...ids]).models;

  return result;
}
