/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Model } from 'ember-cli-mirage';
import { Ref, Workspace } from 'waypoint-pb';

export default Model.extend({
  toProtobuf(): Workspace {
    let result = new Workspace();

    // TODO: result.setActiveTime
    result.setName(this.name);
    // TODO: result.setProjectsList

    return result;
  },

  toProtobufRef(): Ref.Workspace {
    let result = new Ref.Workspace();

    result.setWorkspace(this.name);

    return result;
  },
});
