/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Model, belongsTo, hasMany } from 'miragejs';
import { StatusReport, ResourceCategoryDisplayHint } from 'waypoint-pb';
import { dateToTimestamp } from '../utils';

export default Model.extend({
  statusReport: belongsTo(),
  parent: belongsTo('resource', { inverse: 'children' }),
  children: hasMany('resource', { inverse: 'parent' }),
  declaredResource: belongsTo('resource'),

  toProtobuf(): StatusReport.Resource {
    let result = new StatusReport.Resource();

    result.setCategoryDisplayHint(this.categoryDisplayHintForProtobuf());
    result.setCreatedTime(this.createdTime && dateToTimestamp(this.createdTime));
    result.setDeclaredResource(this.declaredResource?.toProtobufRef());
    result.setHealth(this.healthForProtobuf());
    result.setHealthMessage(this.healthMessage);
    result.setId(this.id);
    result.setName(this.name);
    result.setParentResourceId(this.parent?.id);
    result.setPlatform(this.platform);
    result.setPlatformUrl(this.platformUrl);
    result.setStateJson(JSON.stringify(this.state));
    result.setType(this.type);

    return result;
  },

  categoryDisplayHintForProtobuf(): ResourceCategoryDisplayHint {
    return ResourceCategoryDisplayHint[this.categoryDisplayHint as keyof ResourceCategoryDisplayHint];
  },

  healthForProtobuf(): StatusReport.Resource.Health {
    return StatusReport.Resource.Health[this.health as keyof StatusReport.Resource.Health];
  },
});
