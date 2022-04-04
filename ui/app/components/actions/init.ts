import { Job, Project, QueueJobRequest, Ref } from 'waypoint-pb';

import ApiService from 'waypoint/services/api';
import Component from '@glimmer/component';
import { perform } from 'ember-concurrency-ts';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';

interface ActionsInitArgs {
  project: Project.AsObject;
}

export default class ActionsUp extends Component<ActionsInitArgs> {
  @service api!: ApiService;


  @task
  async init(projectObject: Project.AsObject): Promise<void> {

    let job = new Job();
    let initOp = new Job.InitOp();
    let request = new QueueJobRequest();
    let metadata = this.api.WithMeta();
    let project = new Ref.Project();
    let workspace = new Ref.Workspace();
    let runner = new Ref.Runner();

    workspace.setWorkspace('default');

    project.setProject(projectObject.name);

    // TODO: Figure out how to assign
    // the project to the initOp Job

    runner.setAny(new Ref.RunnerAny());


    job.setInit(initOp);
    job.setTargetRunner(runner);
    job.setWorkspace(workspace);

    request.setJob(job);

    // let response: QueueJobResponse = await this.api.client.queueJob(request, metadata);

    // await perform(this.streamJob, response.getJobId());

    // this.currentMessage = undefined;
    // this.pollModel.route?.refresh();
  }
}
