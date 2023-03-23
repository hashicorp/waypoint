/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Model, belongsTo } from 'miragejs';
import { PushedArtifact } from 'waypoint-pb';

export default Model.extend({
  application: belongsTo(),
  build: belongsTo({ inverse: 'pushedArtifact' }),
  component: belongsTo({ inverse: 'owner' }),
  status: belongsTo({ inverse: 'owner' }),
  workspace: belongsTo(),
  artifact: belongsTo(),

  toProtobuf(): PushedArtifact {
    let result = new PushedArtifact();

    result.setApplication(this.application?.toProtobufRef());
    result.setArtifact(this.artifact?.toProtobuf());
    result.setBuild(this.build?.toProtobuf());
    result.setBuildId(this.build?.id);
    result.setComponent(this.component?.toProtobuf());
    result.setId(this.id);
    result.setJobId(this.jobId);
    result.setSequence(this.sequence);
    result.setStatus(this.status?.toProtobuf());
    result.setTemplateData(this.templateData);
    result.setWorkspace(this.workspace?.toProtobufRef());

    return result;
  },
});
