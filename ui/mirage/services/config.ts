/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { ConfigGetRequest, ConfigGetResponse, ConfigSetRequest, ConfigSetResponse } from 'waypoint-pb';
import { Request, Response } from 'miragejs';
import { RouteHandler } from '../types';
import { decode } from '../helpers/protobufs';

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function get(this: RouteHandler, schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(ConfigGetRequest, requestBody);
  let projectName = requestMsg.getProject()?.getProject();
  let project = schema.projects.findBy({ name: projectName });
  let variables = schema.configVariables.where({ projectId: project.id }).models;
  // The API returns config variables sorted alphabetically by name
  variables.sort((a, b) => (a.name < b.name ? -1 : a.name > b.name ? 1 : 0));
  let variablesList = variables.map((m) => m.toProtobuf());
  let response = new ConfigGetResponse();

  response.setVariablesList(variablesList);

  return this.serialize(response, 'application');
}

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function set(this: RouteHandler, schema: any, { requestBody }: Request): Response {
  // This implementation faithfully recreates the behavior that leads to
  // https://github.com/hashicorp/waypoint/issues/2339.
  // If core changes, we should update this implementation too.

  let requestMsg = decode(ConfigSetRequest, requestBody);
  let vars = requestMsg.toObject().variablesList;

  vars.forEach((attrs) => {
    let project = schema.projects.findBy({ name: attrs.project?.project });
    let configVar = schema.configVariables.findOrCreateBy({ projectId: project.id, name: attrs.name });

    if (attrs.unset !== undefined) {
      configVar.destroy();
    } else {
      configVar.update({
        name: attrs.name,
        pb_static: attrs.pb_static,
        dynamic: attrs.dynamic,
        internal: attrs.internal,
        nameIsPath: attrs.nameIsPath,
      });
    }
  });

  let response = new ConfigSetResponse();

  return this.serialize(response, 'application');
}
