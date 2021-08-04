import Ember from 'ember';
import { logRequestConsole } from './utils';
import failUnhandledRequest from './helpers/fail-unhandled-request';
import { Server } from 'miragejs';

import * as build from './services/build';
import * as project from './services/project';
import * as deployment from './services/deployment';
import * as token from './services/token';
import * as inviteToken from './services/invite-token';
import * as release from './services/release';
import * as versionInfo from './services/version-info';
import * as statusReport from './services/status-report';
import * as job from './services/job';
import * as log from './services/log';
import * as pushedArtifact from './services/pushed-artifact';

export default function (this: Server): void {
  this.namespace = 'hashicorp.waypoint.Waypoint';
  this.urlPrefix = '/grpc';
  this.timing = 0;

  this.pretender.prepareHeaders = (headers) => {
    headers['Content-Type'] = 'application/grpc-web-text';
    headers['X-Grpc-Web'] = '1';
    return headers;
  };

  this.pretender.handledRequest = logRequestConsole;

  if (Ember.testing) {
    this.pretender.unhandledRequest = failUnhandledRequest;
  }

  this.post('/ListBuilds', build.list);
  this.post('/GetBuild', build.get);
  this.post('/ListDeployments', deployment.list);
  this.post('/GetDeployment', deployment.get);
  this.post('/UpsertProject', project.update);
  this.post('/ListProjects', project.list);
  this.post('/GetProject', project.get);
  this.post('/ConvertInviteToken', token.create);
  this.post('/GenerateInviteToken', inviteToken.create);
  this.post('/GenerateLoginToken', token.create);
  this.post('/ListReleases', release.list);
  this.post('/GetRelease', release.get);
  this.post('/GetVersionInfo', versionInfo.get);
  this.post('/ListStatusReports', statusReport.list);
  this.post('/GetLatestStatusReport', statusReport.getLatest);
  this.post('/GetJobStream', job.stream);
  this.post('/GetLogStream', log.stream);
  this.post('/ListPushedArtifacts', pushedArtifact.list);
  this.post('/ExpediteStatusReport', statusReport.expediteStatusReport);

  if (!Ember.testing) {
    // Pass through all other requests
    this.passthrough();
  }
}
