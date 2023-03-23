/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Request, Response } from 'miragejs';
import { RouteHandler } from '../types';
import {
  ListStatusReportsRequest,
  ListStatusReportsResponse,
  ExpediteStatusReportResponse,
} from 'waypoint-pb';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { decode } from '../helpers/protobufs';

// eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/explicit-module-boundary-types
export function list(this: RouteHandler, schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(ListStatusReportsRequest, requestBody);
  let projectName = requestMsg.getApplication()?.getProject();
  let appName = requestMsg.getApplication()?.getApplication();
  let project = schema.projects.findBy({ name: projectName });
  let app = schema.applications.findBy({ projectId: project.id, name: appName });
  // TODO: account for filters
  // TODO: account for workspace
  let statusReports = app.statusReports.models;
  let statusReportProtobufs = statusReports.map((s) => s.toProtobuf());
  let result = new ListStatusReportsResponse();

  result.setStatusReportsList(statusReportProtobufs);

  return this.serialize(result, 'application');
}

export function getLatest(this: RouteHandler): Response {
  return this.serialize(new Empty(), 'application');
}

export function expediteStatusReport(this: RouteHandler): Response {
  // while this is not being used in the current implementation to generate the mocked job id response
  // i'm leaving this here in case we want to update this to handle specific requests
  // let requestMsg = decode(ExpediteStatusReportRequest, requestBody);
  // let ref = requestMsg.getRef();
  let result = new ExpediteStatusReportResponse();
  result.setJobId('JOB_ID');
  return this.serialize(result, 'application');
}
