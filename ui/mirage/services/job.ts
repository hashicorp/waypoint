import { QueueJobRequest, QueueJobResponse } from 'waypoint-pb';
import { Request, Response } from 'miragejs';
import { decode } from '../helpers/protobufs';

export function queue(schema: any, { requestBody }: Request): Response {
  let requestMsg = decode(QueueJobRequest, requestBody);
  let job = requestMsg.getJob();
  let projectName = job.getApplication().getProject();
  let applicationName = job.getApplication().getApplication();
  let workspaceName = job.getWorkspace().getWorkspace();
  let project = schema.projects.findBy({ name: projectName });
  let application = schema.applications.findBy({ name: applicationName, projectId: project.id });
  let workspace = schema.workspaces.findBy({ name: workspaceName });
  let result = new QueueJobResponse();
  let up = job.getUp().toObject();
  let targetRunner = job.getTargetRunner().getId()?.getId() || 'any';
  // TODO: variables?
  // TODO: labels?
  let model = schema.jobs.create({ application, workspace, up, targetRunner });

  result.setJobId(model.id);

  return this.serialize(result, 'application');
}
