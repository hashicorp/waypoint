/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory, trait, association } from 'ember-cli-mirage';
import { fakeId } from '../utils';

export default Factory.extend({
  id: () => fakeId(),
  sequence: (i) => i + 1,

  /**
   * `gitCommitRef` is a shorthand for setting the commit SHA for this build.
   * Provide a value here and it will appear in the `preload` field for the
   * build, wrapped in a `Job.DataSource.Ref` protobuf containing a
   * `Job.Git.Ref` protobuf.
   *
   * @example
   * let build = server.create('build', { gitCommitRef: 'abc123' });
   * console.log(build.toProtobuf().toObject());
   * {
   *   preload: {
   *     jobDataSourceRef: {
   *       git: {
   *         commit: 'abc123'
   *       }
   *     }
   *   }
   * }
   *
   * @type {string | undefined}
   */
  gitCommitRef: undefined,

  afterCreate(build, server) {
    if (!build.workspace) {
      let workspace =
        server.schema.workspaces.findBy({ name: 'default' }) || server.create('workspace', 'default');
      build.update('workspace', workspace);
    }

    build.pushedArtifact?.update('application', build.application);
    build.pushedArtifact?.update('workspace', build.workspace);
  },

  random: trait({
    labels: () => ({
      'common/vcs-ref': '0d56a9f8456b088dd0e4a7b689b842876fd47352',
      'common/vcs-ref-path': 'https://github.com/hashicorp/waypoint/commit/',
    }),
    gitCommitRef: '0d56a9f8456b088dd0e4a7b689b842876fd47352',
    component: association('builder', 'with-random-name'),
    status: association('random'),
    pushedArtifact: association('random'),
  }),

  docker: trait({
    component: association('builder', 'docker'),
  }),

  pack: trait({
    component: association('builder', 'pack'),
  }),

  'seconds-old-success': trait({
    status: association('random', 'success', 'seconds-old'),
  }),

  'minutes-old-success': trait({
    status: association('random', 'success', 'minutes-old'),
  }),

  'hours-old-success': trait({
    status: association('random', 'success', 'hours-old'),
  }),

  'days-old-success': trait({
    status: association('random', 'success', 'days-old'),
  }),
});
