/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Model, belongsTo } from 'ember-cli-mirage';
import { Status } from 'waypoint-pb';
import { dateToTimestamp } from '../utils';

export default Model.extend({
  owner: belongsTo({ polymorphic: true }),

  toProtobuf(): Status {
    let result = new Status();

    result.setCompleteTime(dateToTimestamp(this.completeTime));
    result.setDetails(this.details);
    // result.setError
    // result.setExtension
    result.setStartTime(dateToTimestamp(this.startTime));
    result.setState(Status.State[this.state as keyof typeof Status.State]);

    return result;
  },
});
