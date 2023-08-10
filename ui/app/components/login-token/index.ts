/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import RouterService from '@ember/routing/router-service';
import SessionService from 'ember-simple-auth/services/session';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';

export default class LoginForm extends Component {
  @service session!: SessionService;
  @service router!: RouterService;

  token = '';

  @action
  async login(event?: Event): Promise<void> {
    event?.preventDefault();

    await this.session.authenticate('authenticator:token', this.token);
    this.router.transitionTo('workspaces');
  }
}
