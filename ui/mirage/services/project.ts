import { ListProjectsResponse, GetProjectResponse, UpsertProjectResponse } from 'waypoint-pb';
import { decode } from '../helpers/protobufs';
import { GetProjectRequest, UpsertProjectRequest } from 'waypoint-pb';
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

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function update(schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(UpsertProjectRequest, requestBody);
  let name = requestMsg.getProject().getName();
  let variablesList = requestMsg
    .getProject()
    .getVariablesList()
    .map((v) => v.toObject());
  let model = schema.projects.findBy({ name });

  model.variables = variablesList.map((v) => model.newVariable(v));
  model.save();

  let project = model?.toProtobuf();
  let resp = new UpsertProjectResponse();
  resp.setProject(project);

  return this.serialize(resp, 'application');
}
