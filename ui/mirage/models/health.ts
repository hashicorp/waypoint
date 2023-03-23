/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Model, belongsTo } from 'miragejs';
import { StatusReport } from 'waypoint-pb';

export default Model.extend({
  statusReport: belongsTo({ inverse: 'health' }),
  statusReportList: belongsTo('status-report', { inverse: 'resourcesHealthList' }),

  toProtobuf(): StatusReport.Health {
    let result = new StatusReport.Health();

    result.setHealthStatus(this.status);
    result.setHealthMessage(this.message);

    return result;
  },
});
