/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { GetReleaseRequest, Release, Ref } from 'waypoint-pb';
import ApiService from 'waypoint/services/api';

type Params = { release_id: string };
type Model = Release.AsObject;

export default class ReleaseIdDetail extends Route {
  @service api!: ApiService;

  async model(params: Params): Promise<Model> {
    let req = new GetReleaseRequest();
    let ref = new Ref.Operation();

    ref.setId(params.release_id);
    req.setRef(ref);

    let release: Release = await this.api.client.getRelease(req, this.api.WithMeta());

    return release.toObject();
  }

  redirect(model: Model): void {
    this.transitionTo('workspace.projects.project.app.release', model.sequence);
  }
}
