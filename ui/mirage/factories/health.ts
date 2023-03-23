/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory, trait } from 'ember-cli-mirage';

export default Factory.extend({
  unknown: trait({
    status: 'UNKNOWN',
    message: 'Check is not responding',
    name: 'http',
  }),

  alive: trait({
    status: 'ALIVE',
    message: 'Some resources are alive, application isn’t responding yet',
    name: 'http',
  }),

  ready: trait({
    status: 'READY',
    message: 'All resources are alive, application is responding',
    name: 'http',
  }),

  down: trait({
    status: 'DOWN',
    message: 'No resources are alive, application is not responding',
    name: 'http',
  }),

  partial: trait({
    status: 'PARTIAL',
    message: 'Some resources are down, application is responding',
    name: 'http',
  }),
});
