/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory, trait, association } from 'ember-cli-mirage';

export default Factory.extend({
  'waypoint-ssh': trait({
    url: 'git@github.com:hashicorp/waypoint.git',
    ref: 'main',
    path: 'website',
    ssh: association('example'),
  }),
});
