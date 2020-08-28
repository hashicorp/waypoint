import { Build, ListBuildsResponse, Component, Status, Ref } from 'waypoint-pb';
import { fakeId, fakeComponentForKind, statusRandom, sequenceRandom } from '../utils';

export function createBuild(): Build {
  let build = new Build();
  build.setId(fakeId());

  build.setSequence(sequenceRandom());

  // todo(pearkes): create util
  let workspace = new Ref.Workspace();
  workspace.setWorkspace('default');

  let component = new Component();
  component.setType(Component.Type.BUILDER);
  component.setName(fakeComponentForKind(Component.Type.BUILDER));

  build.setComponent(component);
  build.setStatus(statusRandom());
  build.setWorkspace(workspace);

  build.getLabelsMap().set('common/vcs-ref', '0d56a9f8456b088dd0e4a7b689b842876fd47352');
  build.getLabelsMap().set('common/vcs-ref-path', 'https://github.com/hashicorp/waypoint/commit/');

  return build;
}

export function list(schema: any, { params, requestHeaders }) {
  let resp = new ListBuildsResponse();
  let builds = new Array(createBuild(), createBuild(), createBuild());
  resp.setBuildsList(builds);
  return this.serialize(resp, 'application');
}

export function get(schema: any, { params, requestHeaders }) {
  return this.serialize(createBuild(), 'application');
}
