/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import SessionService from 'ember-simple-auth/services/session';
import Transition from '@ember/routing/-private/transition';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';

const ErrsInvalidToken = ['invalid authentication token', 'Authorization token is not supplied'];

interface ApiError extends Error {
  code: number;
}

export default class Application extends Route {
  @service session!: SessionService;

  async beforeModel(transition: Transition): Promise<void> {
    await this.session.setup();
    await super.beforeModel(transition);
    if (!this.session.isAuthenticated && !transition.to.name.startsWith('auth')) {
      this.session.attemptedTransition = transition;
      this.transitionTo('auth');
    }
  }

  @action
  error(error: ApiError): boolean | void {
    console.log(error);
    let hasAuthError = ErrsInvalidToken.some((msg) => error.message.includes(msg)) || error.code === 16;
    if (hasAuthError) {
      this.session.invalidate();
      this.transitionTo('auth');
    }
    return true;
  }
}
