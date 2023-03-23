/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Model, belongsTo } from 'ember-cli-mirage';
import { Project } from 'waypoint-pb';

export default Model.extend({
  project: belongsTo({ inverse: 'dataSourcePoll' }),

  toProtobuf(): Project.Poll {
    let result = new Project.Poll();

    result.setEnabled(this.enabled ?? false);
    result.setInterval(this.interval ?? '2m');

    return result;
  },
});
