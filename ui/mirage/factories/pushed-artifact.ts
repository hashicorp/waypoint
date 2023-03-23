/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory, association, trait } from 'ember-cli-mirage';
import { fakeId } from '../utils';

export default Factory.extend({
  id: () => fakeId(),
  sequence: (i) => i + 1,

  afterCreate(pushedArtifact, server) {
    if (!pushedArtifact.workspace) {
      let workspace =
        server.schema.workspaces.findBy({ name: 'default' }) || server.create('workspace', 'default');
      pushedArtifact.update('workspace', workspace);
    }
  },

  random: trait({
    component: association('registry', 'with-random-name'),
    status: association('random'),
  }),

  docker: trait({
    component: association('registry', 'docker'),
  }),

  'aws-ecr': trait({
    component: association('registry', 'aws-ecr'),
  }),

  'seconds-old-success': trait({
    status: association('random', 'success', 'seconds-old'),
  }),

  'seconds-old-error': trait({
    status: association('random', 'error', 'seconds-old'),
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
