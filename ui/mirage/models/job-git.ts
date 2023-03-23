/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Model, belongsTo } from 'ember-cli-mirage';
import { Job } from 'waypoint-pb';

export default Model.extend({
  dataSource: belongsTo('job-data-source', { inverse: 'git' }),
  basic: belongsTo('job-git-basic', { inverse: 'parent' }),
  ssh: belongsTo('job-git-ssh', { inverse: 'parent' }),

  toProtobuf(): Job.Git {
    let result = new Job.Git();

    result.setUrl(this.url);
    result.setRef(this.ref);
    result.setPath(this.path);
    result.setIgnoreChangesOutsidePath(this.ignoreChangesOutsidePath ?? true);
    result.setBasic(this.basic?.toProtobuf());
    result.setSsh(this.ssh?.toProtobuf());

    return result;
  },
});
