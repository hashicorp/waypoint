/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import { task, TaskGenerator } from 'ember-concurrency';
import { taskFor } from 'ember-concurrency-ts';
import { Workspace } from 'waypoint-pb';
import ApiService from 'waypoint/services/api';

type Args = {
  current?: string;
  route?: string;
  models?: string[];
  isOpen?: boolean;
};

type Result = { workspaces: Workspace.AsObject[] } | { error: Error };

export default class extends Component<Args> {
  @service api!: ApiService;

  constructor(owner: unknown, args: Args) {
    super(owner, args);

    taskFor(this.loadWorkspaces).perform();
  }

  @task({ restartable: true })
  *loadWorkspaces(): TaskGenerator<Result> {
    try {
      let workspaces = yield this.api.listWorkspaces();
      return { workspaces };
    } catch (error) {
      return { error };
    }
  }

  get shouldExist(): boolean {
    let task = taskFor(this.loadWorkspaces);

    if (task.isRunning) {
      return false;
    }

    let value = task.last?.value;

    if (!value) {
      return false;
    }

    if ('error' in value) {
      return true;
    }

    let { workspaces } = value;

    if (workspaces.length > 1) {
      return true;
    }

    if (workspaces.length === 0) {
      return false;
    }

    if (workspaces[0].name === this.args.current) {
      return false;
    }

    return true;
  }
}
