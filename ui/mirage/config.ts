import { logRequestConsole } from './utils';
import { Server } from 'miragejs';

import * as build from './services/build';
import * as project from './services/project';
import * as deployment from './services/deployment';
import * as token from './services/token';

export default function (this: Server) {
  this.namespace = 'hashicorp.waypoint.Waypoint';
  this.urlPrefix = 'http://localhost:1235';
  this.timing = 0;

  this.pretender.prepareHeaders = (headers) => {
    headers['Content-Type'] = 'application/grpc-web-text';
    headers['X-Grpc-Web'] = '1';
    return headers;
  };

  this.pretender.handledRequest = logRequestConsole;

  this.post('/ListBuilds', build.list);
  this.post('/ListDeployments', deployment.list);
  this.post('/GetDeployment', deployment.get);
  this.post('/ListProjects', project.list);
  this.post('/GetProject', project.get);
  this.post('/ConvertInviteToken', token.create);

  // Pass through all other requests
  this.passthrough();
}
