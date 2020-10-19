import { logRequestConsole } from './utils';
import { Server } from 'miragejs';

import * as build from './services/build';
import * as project from './services/project';
import * as deployment from './services/deployment';
import * as token from './services/token';
import * as inviteToken from './services/invite-token';
import * as release from './services/release';

export default function (this: Server) {
  this.namespace = 'hashicorp.waypoint.Waypoint';
  this.urlPrefix = '/grpc';
  this.timing = 0;

  this.pretender.prepareHeaders = (headers) => {
    headers['Content-Type'] = 'application/grpc-web-text';
    headers['X-Grpc-Web'] = '1';
    return headers;
  };

  this.pretender.handledRequest = logRequestConsole;

  this.post('/ListBuilds', build.list);
  this.post('/GetBuild', build.get);
  this.post('/ListDeployments', deployment.list);
  this.post('/GetDeployment', deployment.get);
  this.post('/ListProjects', project.list);
  this.post('/GetProject', project.get);
  this.post('/ConvertInviteToken', token.create);
  this.post('/GenerateInviteToken', inviteToken.create);
  this.post('/GenerateLoginToken', token.create);
  this.post('/ListReleases', release.list);
  this.post('/GetRelease', release.get);

  // Pass through all other requests
  this.passthrough();
}
