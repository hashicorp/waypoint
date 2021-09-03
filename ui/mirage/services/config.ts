import { ConfigGetRequest, ConfigGetResponse, ConfigSetRequest, ConfigSetResponse } from 'waypoint-pb';
import { Request, Response } from 'miragejs';
import { decode } from '../helpers/protobufs';

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function get(schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(ConfigGetRequest, requestBody);
  let project = requestMsg.getProject();
  let variables = schema.configVariables.where((v) => v.project.name === project.getProject());
  let variablesList = variables.models.toProtobuf();
  let response = new ConfigGetResponse();

  response.setVariablesList(variablesList);

  return this.serialize(response, 'application');
}

export function set(schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(ConfigSetRequest, requestBody);
  let vars = requestMsg.toObject().variablesList;
  vars.forEach((v) => {
    let projName = v.project?.project;
    v.project = null;
    let confVar = schema.configVariables.create(v);
    confVar.project = schema.projects.findBy({ name: projName });
  });
  let response = new ConfigSetResponse();

  return this.serialize(response, 'application');
}
