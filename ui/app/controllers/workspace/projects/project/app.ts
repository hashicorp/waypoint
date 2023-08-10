/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import RouterService from '@ember/routing/router-service';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

export default class extends Controller {
  @tracked isSwitchingWorkspace = false;
  @service router!: RouterService;

  /**
   * Returns a suitable “pivot” route when switching between workspaces.
   * For example, if you’re looking at a list of builds for production
   * then it makes sense to switch to the list of builds for staging.
   * However, if you’re looking at an individual deployment in production
   * then it wouldn’t make sense to switch to that same deployment in staging,
   * because it won’t be there. Instead, you want to jump up a level to see
   * the list of all deployments for production.
   */
  get workspaceSwitcherTargetRoute(): RouterService['currentRoute'] {
    let result = this.router.currentRoute;

    while (result.parent && result.parent.localName !== 'app') {
      result = result.parent;
    }

    return result;
  }

  /**
   * Returns the array of models to go with `workspaceSwitcherTargetRoute`,
   * except without the leading workspace model (which the switcher will provide).
   */
  get workspaceSwitcherTargetModels(): unknown[] {
    let result: unknown[] = [];
    let route = this.workspaceSwitcherTargetRoute;

    while (route.parent && route.parent.localName !== 'workspace') {
      result = [...Object.values(route.params), ...result];
      route = route.parent;
    }

    return result;
  }
}
