import { Build, Ref, ListProjectsResponse } from 'waypoint-pb';
import { fakeId } from '../utils';
import faker from '../faker';
import { dasherize } from '@ember/string';

function createProjectRef(): Ref.Project {
  let build = new Build();
  build.setId(fakeId());

  // todo(pearkes): create util
  let workspace = new Ref.Workspace();
  workspace.setWorkspace('default');

  let project = new Ref.Project();
  project.setProject(dasherize(faker.hacker.noun()));

  return project;
}

export function list(schema: any, { params, requestHeaders }) {
  let resp = new ListProjectsResponse();
  let projs = new Array(createProjectRef(), createProjectRef());
  resp.setProjectsList(projs);
  return this.serialize(resp, 'application');
}
