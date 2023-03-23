/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Model, hasMany } from 'miragejs';
import { Generation } from 'waypoint-pb';

export default Model.extend({
  deployment: hasMany(),

  toProtobuf(): Generation {
    let result = new Generation();

    result.setId(this.id);
    result.setInitialSequence(this.initialSequence);

    return result;
  },
});
