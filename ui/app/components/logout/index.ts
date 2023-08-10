/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import RouterService from '@ember/routing/router-service';
import SessionService from 'ember-simple-auth/services/session';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';

export default class Logout extends Component {
  @service session!: SessionService;
  @service router!: RouterService;

  @action
  async logout(): Promise<void> {
    await this.session.invalidate();
    this.router.transitionTo('auth');
  }
}
