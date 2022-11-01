import { ListProjectsResponse, GetProjectResponse, UpsertProjectResponse } from 'waypoint-pb';
import { decode } from '../helpers/protobufs';
import { UI, GetProjectRequest, UpsertProjectRequest, Job, Project } from 'waypoint-pb';
import { Request, Response } from 'miragejs';
import { RouteHandler } from '../types';

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function list(this: RouteHandler, schema: any): Response {
  let resp = new ListProjectsResponse();
  let projectRefs = schema.projects.all().models.map((p) => p.toProtobufRef());

  resp.setProjectsList(projectRefs);

  return this.serialize(resp, 'application');
}

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function get(this: RouteHandler, schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(GetProjectRequest, requestBody);
  let name = requestMsg.getProject()?.getProject();
  let model = schema.projects.findBy({ name });
  let resp = new GetProjectResponse();
  let project = model?.toProtobuf();

  resp.setProject(project);

  return this.serialize(resp, 'application');
}

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function uiGet(this: RouteHandler, schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(UI.GetProjectRequest, requestBody);
  let name = requestMsg.getProject()?.getProject();
  let model = schema.projects.findBy({ name });
  let resp = new UI.GetProjectResponse();
  let project = model?.toProtobuf();

  resp.setProject(project);

  // TODO(jgwhite): sideload latest init job
  if (model.name === 'init-test') {
    let initJob = new Job();
    initJob.setInit(new Job.InitOp());
    initJob.setState(Job.State.RUNNING);
    resp.setLatestInitJob(initJob);
  }

  return this.serialize(resp, 'application');
}

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any
export function upsert(this: RouteHandler, schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(UpsertProjectRequest, requestBody);
  let name = requestMsg.getProject()?.getName();
  let variablesList = requestMsg
    .getProject()
    ?.getVariablesList()
    .map((v) => v.toObject());
  let dataSource = requestMsg.getProject()?.getDataSource();
  let poll = requestMsg.getProject()?.getDataSourcePoll();
  let model = schema.projects.findOrCreateBy({ name });

  model.variables = variablesList?.map((v) => model.newVariable(v));
  model.dataSource = dataSourceFromProto(schema, dataSource);
  model.dataSourcePoll = pollFromProto(schema, poll);
  model.save();

  let project = model?.toProtobuf();
  let resp = new UpsertProjectResponse();
  resp.setProject(project);

  return this.serialize(resp, 'application');
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function dataSourceFromProto(schema: any, dataSource: Job.DataSource | undefined): any {
  if (!dataSource) {
    return;
  }

  let dataSourceObj = dataSource.toObject();
  let result = schema.new('job-data-source');

  if (dataSourceObj.git) {
    let { url, ref, path, ignoreChangesOutsidePath } = dataSourceObj.git;
    let git = result.newGit({ url, ref, path, ignoreChangesOutsidePath });

    if (dataSourceObj.git.ssh) {
      git.newSsh(dataSourceObj.git.ssh);
    }

    if (dataSourceObj.git.basic) {
      git.newBasic(dataSourceObj.git.basic);
    }
  }

  return result;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function pollFromProto(schema: any, poll: Project.Poll | undefined): any {
  if (!poll) {
    return;
  }

  let result = schema.new('project-poll', poll.toObject());

  return result;
}
