/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { ListReleasesRequest, ListReleasesResponse, GetReleaseRequest, UI } from 'waypoint-pb';
import { Request, Response } from 'miragejs';
import { RouteHandler } from '../types';
import { decode } from '../helpers/protobufs';

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function list(this: RouteHandler, schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(ListReleasesRequest, requestBody);
  let projectName = requestMsg.getApplication()?.getProject();
  let appName = requestMsg.getApplication()?.getApplication();
  let workspaceName = requestMsg.getWorkspace()?.getWorkspace();
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
export function ui_list(this: RouteHandler, schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(UI.ListReleasesRequest, requestBody);
  let projectName = requestMsg.getApplication()?.getProject();
  let appName = requestMsg.getApplication()?.getApplication();
  let workspaceName = requestMsg.getWorkspace()?.getWorkspace();
  let project = schema.projects.findBy({ name: projectName });
  let application = schema.applications.findBy({ name: appName, projectId: project?.id });
  let workspace = schema.workspaces.findBy({ name: workspaceName });
  let releases = schema.releases
    .where({ applicationId: application?.id, workspaceId: workspace?.id })
    .models.sort((a, b) => b.sequence - a.sequence);
  let bundles: UI.ReleaseBundle[] = releases.map((release) => {
    let bundle = new UI.ReleaseBundle();

    bundle.setRelease(release.toProtobuf());
    bundle.setLatestStatusReport(release.statusReport?.toProtobuf());

    return bundle;
  });
  let resp = new UI.ListReleasesResponse();

  resp.setReleasesList(bundles);

  return this.serialize(resp, 'application');
}
// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function get(this: RouteHandler, schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(GetReleaseRequest, requestBody);
  let id = requestMsg.getRef()?.getId();
  let model = schema.releases.find(id);
  let protobuf = model?.toProtobuf();

  return this.serialize(protobuf, 'application');
}
