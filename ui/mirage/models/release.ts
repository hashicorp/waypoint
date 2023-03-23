/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Model, belongsTo } from 'ember-cli-mirage';
import { Release, Operation } from 'waypoint-pb';

const { PhysicalState } = Operation;
type StateName = keyof typeof PhysicalState;

export default Model.extend({
  application: belongsTo(),
  workspace: belongsTo(),
  deployment: belongsTo(),
  status: belongsTo({ inverse: 'owner' }),
  component: belongsTo({ inverse: 'owner' }),
  statusReport: belongsTo({ inverse: 'target' }),

  toProtobuf(): Release {
    let result = new Release();

    result.setApplication(this.application?.toProtobufRef());
    result.setComponent(this.component?.toProtobuf());
    result.setDeploymentId(this.deployment?.id);
    result.setId(this.id);
    result.setJobId(this.jobId);
    result.setPreload(this.preloadProtobuf());
    result.setSequence(this.sequence);
    result.setState(PhysicalState[this.state as StateName]);
    result.setStatus(this.status?.toProtobuf());
    result.setTemplateData(this.templateData);
    result.setUrl(this.url);
    result.setWorkspace(this.workspace?.toProtobufRef());

    for (let [key, value] of Object.entries<string>(this.labels ?? {})) {
      result.getLabelsMap().set(key, value);
    }

    return result;
  },

  preloadProtobuf(): Release.Preload {
    let result = new Release.Preload();

    result.setArtifact(this.deployment?.build?.pushedArtifact?.toProtobuf());
    result.setBuild(this.deployment?.build?.toProtobuf());
    result.setDeployment(this.deployment?.toProtobuf());

    return result;
  },
});
