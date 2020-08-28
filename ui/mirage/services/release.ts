import { Component, Ref, Release, ListReleasesResponse } from 'waypoint-pb';
import { fakeId, fakeComponentForKind, statusRandom, sequenceRandom } from '../utils';
import { createDeployment } from './deployment';

function createRelease(): Release {
  let release = new Release();
  release.setId(fakeId());

  // todo(pearkes): create util
  let workspace = new Ref.Workspace();
  workspace.setWorkspace('default');

  let component = new Component();
  component.setType(Component.Type.RELEASEMANAGER);
  component.setName(fakeComponentForKind(Component.Type.RELEASEMANAGER));

  release.setSequence(sequenceRandom());
  release.setStatus(statusRandom());
  release.setComponent(component);
  release.setWorkspace(workspace);
  release.setDeploymentId(createDeployment().getId());

  release.setUrl(`https://wildly-intent-honeybee--${release.getId().toLowerCase()}.alpha.waypoint.run`);

  release.setState(3);

  release.getLabelsMap().set('common/vcs-ref', '0d56a9f8456b088dd0e4a7b689b842876fd47352');
  release.getLabelsMap().set('common/vcs-ref-path', 'https://github.com/hashicorp/waypoint/commit/');

  return release;
}

export function list(schema: any, { params, requestHeaders }) {
  let resp = new ListReleasesResponse();
  let releases = new Array(createRelease(), createRelease(), createRelease(), createRelease());
  resp.setReleasesList(releases);
  return this.serialize(resp, 'application');
}

export function get(schema: any, { params, requestHeaders }) {
  let deploy = createRelease();
  return this.serialize(deploy, 'application');
}
