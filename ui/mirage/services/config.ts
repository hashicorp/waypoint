import { ConfigGetRequest, ConfigSetRequest, ConfigSetResponse } from 'waypoint-pb';
import { Request, Response } from 'miragejs';
import { decode } from '../helpers/protobufs';
import configVariable from '../factories/config-variable';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function get(schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(ConfigGetRequest, requestBody);
  let project = requestMsg.getProject();
  let model = schema.configVariables.where((v) => v.project.name === project.getProject());
  let variablesList = model?.toProtobuf();

  return this.serialize(variablesList, 'application');
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
  // This API endpoint returns an empty {} response, not even Empty(), so serialization somehow fails and the code below doesn't work.
  // let resp = new ConfigSetResponse();
  // return this.serialize(resp, 'application');
  return {};
}
