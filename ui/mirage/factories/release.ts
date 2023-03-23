/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory, trait, association } from 'ember-cli-mirage';
import { fakeId } from '../utils';

export default Factory.extend({
  id: () => fakeId(),
  sequence: (i) => i + 1,

  afterCreate(release, server) {
    if (!release.workspace) {
      let workspace =
        server.schema.workspaces.findBy({ name: 'default' }) || server.create('workspace', 'default');
      release.update('workspace', workspace);
    }
  },

  random: trait({
    component: association('release-manager', 'with-random-name'),
    status: association('random'),
    state: 'CREATED',
    labels: () => ({
      'common/vcs-ref': '0d56a9f8456b088dd0e4a7b689b842876fd47352',
      'common/vcs-ref-path': 'https://github.com/hashicorp/waypoint/commit/',
    }),
    url: 'https://wp-matrix.example',
  }),

  nomad: trait({
    component: association('release-manager', { name: 'nomad' }),
  }),

  docker: trait({
    component: association('release-manager', { name: 'docker' }),
  }),

  'nomad-jobspec': trait({
    component: association('release-manager', { name: 'nomad-jobspec' }),
  }),

  'nomad-jobspec-canary': trait({
    component: association('release-manager', { name: 'nomad-jobspec-canary' }),
  }),

  kubernetes: trait({
    component: association('release-manager', { name: 'kubernetes' }),
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
