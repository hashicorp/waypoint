/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Model, belongsTo } from 'ember-cli-mirage';
import { Job } from 'waypoint-pb';

export default Model.extend({
  parent: belongsTo('job-git', { inverse: 'basic' }),

  toProtobuf(): Job.Git.Basic {
    let result = new Job.Git.Basic();

    result.setUsername(this.username);
    result.setPassword(this.password);

    return result;
  },
});
