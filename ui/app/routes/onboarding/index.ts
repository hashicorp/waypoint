/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';

export default class OnboardingIndex extends Route {
  redirect(): void {
    this.transitionTo('onboarding.install');
  }
}
