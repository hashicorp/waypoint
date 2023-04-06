/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import '@ember/routing/router-service';
import Route from '@ember/routing/route';

declare module '@ember/routing/router-service' {
  type Transition = ReturnType<Route['transitionTo']>;

  export default interface RouterService {
    // This method comes from ember-router-service-refresh-polyfill,
    // which does not provide its own type declarations.
    refresh(pivotRouteName?: string): Transition;
  }
}
