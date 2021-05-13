import Controller from '@ember/controller';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import ApiService from 'waypoint/services/api';
import { Ref, Job, QueueJobRequest } from 'waypoint-pb';

export default class extends Controller {
  @service api!: ApiService;

  @action
  async up(appObject: Ref.Application.AsObject): Promise<void> {
    let job = new Job();
    let upOp = new Job.UpOp();
    let request = new QueueJobRequest();
    let metadata = this.api.WithMeta();
    let application = new Ref.Application();
    let workspace = new Ref.Workspace();
    let runner = new Ref.Runner();

    if (!appObject.project) {
      return;
    }

    workspace.setWorkspace('default');

    application.setApplication(appObject.application);
    application.setProject(appObject.project);

    runner.setAny(new Ref.RunnerAny());

    job.setApplication(application);
    job.setUp(upOp);
    job.setTargetRunner(runner);
    job.setWorkspace(workspace);

    request.setJob(job);

    try {
      let response = await this.api.client.queueJob(request, metadata);
      console.log(response);
    } catch (error) {
      console.error(error);
    }
  }
}
