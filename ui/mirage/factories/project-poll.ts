/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory, trait } from 'ember-cli-mirage';

export default Factory.extend({
  'every-2-minutes': trait({
    enabled: true,
    interval: '2m',
  }),
});
