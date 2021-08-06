import { ConfigGetRequest } from 'waypoint-pb';
import { Request, Response } from 'miragejs';
import { decode } from '../helpers/protobufs';
import configVariable from '../factories/config-variable';

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function get(schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(ConfigGetRequest, requestBody);
  // let project = requestMsg.getProject();
  let model = schema.projects.findAll(configVariable);
  let variablesList = model?.toProtobuf();

  return this.serialize(variablesList, 'application');
}
