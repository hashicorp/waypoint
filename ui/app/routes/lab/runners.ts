/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { ListRunnersRequest, Runner } from 'waypoint-pb';

type Model = Runner.AsObject[];

export default class extends Route {
  @service api!: ApiService;

  async model(): Promise<Model> {
    let request = new ListRunnersRequest();
    let response = await this.api.client.listRunners(request, this.api.WithMeta());

    return response.toObject().runnersList;
  }
}
