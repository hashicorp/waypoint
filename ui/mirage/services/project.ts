import { ListProjectsResponse, GetProjectResponse } from 'waypoint-pb';
import { decode } from '../helpers/protobufs';
import { GetProjectRequest } from 'waypoint-pb';
import { Request, Response } from 'miragejs';

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function list(schema: any): Response {
  let resp = new ListProjectsResponse();
  let projectRefs = schema.projects.all().models.map((p) => p.toProtobufRef());

  resp.setProjectsList(projectRefs);

  return this.serialize(resp, 'application');
}

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function get(schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(GetProjectRequest, requestBody);
  let name = requestMsg.getProject().getProject();
  let model = schema.projects.findBy({ name });
  let resp = new GetProjectResponse();
  let project = model?.toProtobuf();

  resp.setProject(project);

  return this.serialize(resp, 'application');
}
