/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory, trait } from 'ember-cli-mirage';

import faker from '../faker';

export default Factory.extend({
  name: () => faker.company.companyName(),
  displayName: () => faker.company.companyName(),
  kind: 0,
  google: trait({
    name: 'google',
    displayName: 'google',
  }),
});
