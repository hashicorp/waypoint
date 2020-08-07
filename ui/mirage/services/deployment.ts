import { Component, Status, Ref, Deployment, ListDeploymentsResponse } from 'waypoint-pb';
import { fakeId, fakeComponentForKind } from '../utils';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { subMinutes } from 'date-fns';

function createDeployment(): Deployment {
  let deploy = new Deployment();
  deploy.setId(fakeId());

  // todo(pearkes): create util
  let workspace = new Ref.Workspace();
  workspace.setWorkspace('default');

  let component = new Component();
  component.setType(Component.Type.REGISTRY);
  component.setName(fakeComponentForKind(Component.Type.REGISTRY));

  // todo(pearkes): random state
  let status = new Status();
  status.setState(Status.State.SUCCESS);

  // todo(pearkes): helpers
  let timestamp = new Timestamp();
  let result = Math.floor(subMinutes(new Date(), 30).getTime() / 1000);
  timestamp.setSeconds(result);

  // Same thing for now
  status.setCompleteTime(timestamp);
  status.setStartTime(timestamp);

  deploy.setComponent(component);
  deploy.setStatus(status);
  deploy.setWorkspace(workspace);

  return deploy;
}

export function list(schema: any, { params, requestHeaders }) {
  let resp = new ListDeploymentsResponse();
  let deploys = new Array(createDeployment(), createDeployment(), createDeployment(), createDeployment());
  resp.setDeploymentsList(deploys);
  return this.serialize(resp, 'application');
}

export function get(schema: any, { params, requestHeaders }) {
  let deploy = createDeployment();
  return this.serialize(deploy, 'application');
}
