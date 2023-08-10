/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { Build, GetBuildRequest, Ref } from 'waypoint-pb';
import ApiService from 'waypoint/services/api';

type Params = { build_id: string };
type Model = Build.AsObject;

export default class WorkspaceProjectsProjectAppBuildId extends Route {
  @service api!: ApiService;

  async model(params: Params): Promise<Model> {
    let req = new GetBuildRequest();
    let ref = new Ref.Operation();

    ref.setId(params.build_id);
    req.setRef(ref);

    let build = await this.api.client.getBuild(req, this.api.WithMeta());

    return build.toObject();
  }

  redirect(model: Model): void {
    this.transitionTo('workspace.projects.project.app.build', model.sequence);
  }
}
