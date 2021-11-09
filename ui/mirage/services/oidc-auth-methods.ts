import { Request, Response, RouteHandler } from 'miragejs';

import { ListOIDCAuthMethodsResponse } from 'waypoint-pb';

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function list(this: RouteHandler, schema: any, { requestBody }: Request): Response {
  let authMethods = schema.authMethods.all();
  let authMethodsProtos = authMethods.models?.map((model) => model?.toProtobuf());

  let resp = new ListOIDCAuthMethodsResponse();
  resp.setAuthMethodsList(authMethodsProtos);
  return this.serialize(resp, 'application');
}
