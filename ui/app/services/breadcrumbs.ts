/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { getOwner } from '@ember/application';
import Service, { inject as service } from '@ember/service';
import { computed } from '@ember/object';
import classic from 'ember-classic-decorator';

export interface Breadcrumb {
  label: string;
  route: string;
  icon?: string;
  isPending?: boolean;
}

type BreadcrumbBuilder = (model: unknown) => Breadcrumb[];

@classic
export default class BreadcrumbsService extends Service {
  @service router;

  // currentURL is only used to listen to all transitions.
  // currentRouteName has all information necessary to compute breadcrumbs,
  // but it doesn't change when a transition to the same route with a different
  // model occurs.
  @computed('router.{currentURL,currentRouteName}')
  get breadcrumbs(): Breadcrumb[] {
    let owner = getOwner(this);
    let allRoutes = (this.router.currentRouteName || '')
      .split('.')
      .without('')
      .map((_segment, index, allSegments) => allSegments.slice(0, index + 1).join('.'));

    let crumbs: Breadcrumb[] = [];
    allRoutes.forEach((routeName) => {
      let route = owner.lookup(`route:${routeName}`);

      // Routes can reset the breadcrumb trail to start anew even
      // if the route is deeply nested.
      if (route.resetBreadcrumbs) {
        crumbs = [];
      }

      // Breadcrumbs are either an array of static crumbs
      // or a function that returns breadcrumbs given the current
      // model for the route's controller.
      let breadcrumbs: BreadcrumbBuilder | Breadcrumb[] = route.breadcrumbs || [];

      if (typeof breadcrumbs === 'function') {
        breadcrumbs = breadcrumbs(route.get('controller.model')) || [];
      }

      crumbs.push(...breadcrumbs);
    });

    return crumbs;
  }
}
