import { Component, Ref, Deployment, ListDeploymentsResponse } from 'waypoint-pb';
import { fakeId, fakeComponentForKind, statusRandom, sequenceRandom } from '../utils';
import { createBuild } from './build';

export function createDeployment(): Deployment {
  let deploy = new Deployment();
  deploy.setId(fakeId());

  // todo(pearkes): create util
  let workspace = new Ref.Workspace();
  workspace.setWorkspace('default');

  let component = new Component();
  component.setType(Component.Type.PLATFORM);
  component.setName(fakeComponentForKind(Component.Type.PLATFORM));

  deploy.setSequence(sequenceRandom());
  deploy.setStatus(statusRandom());
  deploy.setComponent(component);
  deploy.setWorkspace(workspace);
  deploy.setBuild(createBuild());

  deploy.setState(3);

  deploy.getLabelsMap().set('common/vcs-ref', '0d56a9f8456b088dd0e4a7b689b842876fd47352');
  deploy.getLabelsMap().set('common/vcs-ref-path', 'https://github.com/hashicorp/waypoint/commit/');

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
