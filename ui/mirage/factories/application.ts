/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory, trait } from 'ember-cli-mirage';
import faker from '../faker';

export default Factory.extend({
  fileChangeSignal: 'HUP',

  simple: trait({
    name: 'simple-application',
  }),

  'with-random-name': trait({
    name: () => `wp-${faker.hacker.noun()}`,
  }),
});
