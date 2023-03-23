/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Model, belongsTo } from 'ember-cli-mirage';
import { Variable } from 'waypoint-pb';

export default Model.extend({
  project: belongsTo(),

  toProtobuf(): Variable {
    let result = new Variable();

    result.setServer();
    result.setName(this.name);
    result.setSensitive(this.sensitive);
    if (this.hcl) {
      result.setStr('');
      result.setHcl(this.hcl);
    } else {
      if (this.str) {
        result.setHcl('');
        result.setStr(this.str);
      }
    }

    return result;
  },
});
