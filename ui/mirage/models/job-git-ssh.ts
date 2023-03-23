/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Model, belongsTo } from 'ember-cli-mirage';
import { Job } from 'waypoint-pb';

export default Model.extend({
  parent: belongsTo('job-git', { inverse: 'ssh' }),

  toProtobuf(): Job.Git.SSH {
    let result = new Job.Git.SSH();

    result.setUser(this.user);
    result.setPassword(this.password);
    result.setPrivateKeyPem(this.privateKeyPem);

    return result;
  },
});
