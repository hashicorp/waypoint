/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory, trait, association } from 'ember-cli-mirage';

export default Factory.extend({
  'marketing-public': trait({
    git: association('waypoint-ssh'),
  }),
});
