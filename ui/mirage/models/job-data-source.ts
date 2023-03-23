/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Model, belongsTo } from 'ember-cli-mirage';
import { Job } from 'waypoint-pb';

export default Model.extend({
  project: belongsTo({ inverse: 'dataSource' }),
  git: belongsTo('job-git', { inverse: 'dataSource' }),

  toProtobuf(): Job.DataSource {
    let result = new Job.DataSource();

    if (this.git) {
      result.setGit(this.git.toProtobuf());
    } else {
      result.setLocal(new Job.Local());
    }

    return result;
  },
});
