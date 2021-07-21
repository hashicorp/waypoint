import { Request, Response } from 'miragejs';
import { ListStatusReportsRequest, ListStatusReportsResponse } from 'waypoint-pb';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { decode } from '../helpers/protobufs';

// eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/explicit-module-boundary-types
export function list(schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(ListStatusReportsRequest, requestBody);
  let projectName = requestMsg.getApplication().getProject();
  let appName = requestMsg.getApplication().getApplication();
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

export function getLatest(): Response {
  return this.serialize(new Empty(), 'application');
}
