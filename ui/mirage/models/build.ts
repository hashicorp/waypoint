/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Model, belongsTo, hasMany } from 'ember-cli-mirage';
import { Build, Job } from 'waypoint-pb';

export default Model.extend({
  application: belongsTo(),
  workspace: belongsTo(),
  component: belongsTo({ inverse: 'owner' }),
  status: belongsTo({ inverse: 'owner' }),
  deployments: hasMany(),
  artifact: belongsTo(),
  pushedArtifact: belongsTo({ inverse: 'build' }),

  toProtobuf(): Build {
    let result = new Build();

    result.setApplication(this.application?.toProtobufRef());
    result.setArtifact(this.artifact?.toProtobuf());
    result.setComponent(this.component?.toProtobuf());
    result.setId(this.id);
    result.setJobId(this.JobId);
    result.setPreload(this.preloadProtobuf());
    result.setSequence(this.sequence);
    result.setStatus(this.status?.toProtobuf());
    result.setTemplateData(this.templateData);
    result.setWorkspace(this.workspace?.toProtobufRef());

    for (let [key, value] of Object.entries<string>(this.labels ?? {})) {
      result.getLabelsMap().set(key, value);
    }

    return result;
  },

  preloadProtobuf(): Build.Preload {
    let result = new Build.Preload();

    if (this.gitCommitRef) {
      let dataSourceRef = new Job.DataSource.Ref();
      let gitRef = new Job.Git.Ref();

      gitRef.setCommit(this.gitCommitRef);
      dataSourceRef.setGit(gitRef);

      result.setJobDataSourceRef(dataSourceRef);
    }

    return result;
  },
});
