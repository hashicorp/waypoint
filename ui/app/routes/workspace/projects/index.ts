/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { Ref } from 'waypoint-pb';
import ApiService from 'waypoint/services/api';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import PollModelService from 'waypoint/services/poll-model';
import ProjectsIndex from 'waypoint/controllers/workspace/projects';

type Model = Ref.Project.AsObject[];

export default class Index extends Route {
  @service api!: ApiService;
  @service pollModel!: PollModelService;

  async model(): Promise<Model> {
    let resp = await this.api.client.listProjects(new Empty(), this.api.WithMeta());
    let projects = resp.getProjectsList().map((p) => p.toObject());

    return projects;
  }

  afterModel(): void {
    this.pollModel.setup(this);
  }

  resetController(controller: ProjectsIndex): void {
    // Clear the CLI parameter when we leave this route or update the model
    controller.set('cli', null);
  }
}
