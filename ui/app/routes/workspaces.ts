/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApiService from 'waypoint/services/api';
import Route from '@ember/routing/route';
import SessionService from 'ember-simple-auth/services/session';
import { Workspace } from 'waypoint-pb';
import { inject as service } from '@ember/service';

export default class Workspaces extends Route {
  @service api!: ApiService;
  @service session!: SessionService;

  async redirect(): Promise<void> {
    let workspaces: Workspace.AsObject[] = [];
    try {
      workspaces = await this.api.listWorkspaces();
    } catch (error) {
      // Send Authentication or other error if it exists
      this.send('error', error);
    }
    let workspaceNames = workspaces.map((w) => w.name).sort();
    let storedWorkspaceName = this.session.data.workspace;

    if (storedWorkspaceName && workspaceNames.includes(storedWorkspaceName)) {
      this.transitionTo('workspace', storedWorkspaceName);
      return;
    }

    if (workspaceNames.includes('default')) {
      this.transitionTo('workspace', 'default');
      return;
    }

    if (workspaceNames.length) {
      this.transitionTo('workspace', workspaceNames[0]);
      return;
    }

    this.transitionTo('workspace', 'default');
  }
}
