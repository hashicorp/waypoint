/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { Ref } from 'waypoint-pb';
import SessionService from 'ember-simple-auth/services/session';
import { inject as service } from '@ember/service';

export type Params = { workspace_id: string };
export type Model = Ref.Workspace.AsObject;

export default class Workspace extends Route {
  @service session!: SessionService;

  async model(params: Params): Promise<Model> {
    // Workspace "id" which is a name, based on URL param
    let ws = new Ref.Workspace();
    ws.setWorkspace(params.workspace_id);

    return ws.toObject();
  }

  afterModel(model: Model): void {
    this.storeWorkspace(model.workspace);
  }

  storeWorkspace(workspace: string): void {
    this.session.set('data.workspace', workspace);
  }
}
