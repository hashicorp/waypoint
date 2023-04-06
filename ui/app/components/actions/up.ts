/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import {
  GetJobRequest,
  GetJobStreamResponse,
  Job,
  QueueJobRequest,
  QueueJobResponse,
  Ref,
} from 'waypoint-pb';

import ApiService from 'waypoint/services/api';
import Component from '@glimmer/component';
import PollModelService from 'waypoint/services/poll-model';
import { perform } from 'ember-concurrency-ts';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';

interface ActionsUpArgs {
  application: Ref.Application.AsObject;
}

export default class ActionsUp extends Component<ActionsUpArgs> {
  @service api!: ApiService;
  @service pollModel!: PollModelService;

  @tracked currentMessage?: string;

  @task
  async up(appObject: Ref.Application.AsObject): Promise<void> {
    this.currentMessage = undefined;

    let job = new Job();
    let upOp = new Job.UpOp();
    let request = new QueueJobRequest();
    let metadata = this.api.WithMeta();
    let application = new Ref.Application();
    let workspace = new Ref.Workspace();
    let runner = new Ref.Runner();

    workspace.setWorkspace('default');

    application.setApplication(appObject.application);
    application.setProject(appObject.project);

    runner.setAny(new Ref.RunnerAny());

    job.setApplication(application);
    job.setUp(upOp);
    job.setTargetRunner(runner);
    job.setWorkspace(workspace);

    request.setJob(job);

    let response: QueueJobResponse = await this.api.client.queueJob(request, metadata);

    await perform(this.streamJob, response.getJobId());

    this.currentMessage = undefined;
    this.pollModel.route?.refresh();
  }

  @task
  async streamJob(jobId: string): Promise<void> {
    let request = new GetJobRequest();
    let metadata = this.api.WithMeta();

    request.setJobId(jobId);

    let stream = this.api.client.getJobStream(request, metadata);

    await new Promise<void>((resolve, reject) => {
      stream.on('status', (resp) => console.log(resp));
      stream.on('metadata', (resp) => console.log(resp));
      stream.on('data', (resp: GetJobStreamResponse) => {
        resp
          .getTerminal()
          ?.getEventsList()
          .map((event) => event.getLine()?.getMsg())
          .filter(notEmpty)
          .forEach((msg) => {
            this.pollModel.route?.refresh();
            this.currentMessage = msg;
          });
      });
      stream.on('error', reject);
      stream.on('end', () => {
        resolve();
      });
    });
  }
}

function notEmpty<TValue>(value: TValue | null | undefined): value is TValue {
  return value !== null && value !== undefined;
}
