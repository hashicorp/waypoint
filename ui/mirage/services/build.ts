import { Build, ListBuildsResponse, Component, Status, Ref } from 'waypoint-pb';
import { fakeId, fakeComponentForKind, statusRandom } from '../utils';

const buildSequence = 0;

function createBuild(): Build {
  let build = new Build();
  build.setId(fakeId());

  build.setSequence(buildSequence + 1);

  // todo(pearkes): create util
  let workspace = new Ref.Workspace();
  workspace.setWorkspace('default');

  let component = new Component();
  component.setType(Component.Type.BUILDER);
  component.setName(fakeComponentForKind(Component.Type.BUILDER));

  build.setComponent(component);
  build.setStatus(statusRandom());
  build.setWorkspace(workspace);

  return build;
}

export function list(schema: any, { params, requestHeaders }) {
  let resp = new ListBuildsResponse();
  let builds = new Array(createBuild(), createBuild(), createBuild(), createBuild());
  resp.setBuildsList(builds);
  return this.serialize(resp, 'application');
}

export function get(schema: any, { params, requestHeaders }) {
  return this.serialize(createBuild(), 'application');
}
