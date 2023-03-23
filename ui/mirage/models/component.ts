/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Model, belongsTo } from 'ember-cli-mirage';
import { Component } from 'waypoint-pb';

// eslint-disable-next-line ember/require-tagless-components
export default Model.extend({
  owner: belongsTo({ polymorphic: true }),

  toProtobuf(): Component {
    let result = new Component();

    result.setName(this.name);
    result.setType(Component.Type[this.type as keyof typeof Component.Type]);

    return result;
  },
});
