import { Build, Ref, ListProjectsResponse, GetProjectResponse, Project, Application } from 'waypoint-pb';
import { fakeId } from '../utils';
import faker from '../faker';
import { dasherize } from '@ember/string';
import { create } from 'domain';

const projectName = 'marketing-public';

function createProjectRef(): Ref.Project {
  let build = new Build();
  build.setId(fakeId());

  // todo(pearkes): create util
  let workspace = new Ref.Workspace();
  workspace.setWorkspace('default');

  let project = new Ref.Project();
  project.setProject(projectName);

  return project;
}

function createApp(): Application {
  let app = new Application();
  app.setName(`wp-${faker.hacker.noun()}`);

  return app;
}

function createProject(): Project {
  let proj = new Project();
  proj.setName(projectName);
  proj.setApplicationsList([createApp()]);

  return proj;
}

export function list(schema: any, { params, requestHeaders }) {
  let resp = new ListProjectsResponse();
  let projs = new Array(createProjectRef());
  resp.setProjectsList(projs);
  return this.serialize(resp, 'application');
}

export function get(schema: any, { params, requestHeaders }) {
  let resp = new GetProjectResponse();
  let proj = createProject();
  resp.setProject(proj);
  return this.serialize(resp, 'application');
}
