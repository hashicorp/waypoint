/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory, trait, association } from 'ember-cli-mirage';

export default Factory.extend({
  afterCreate(statusReport, server) {
    if (!statusReport.workspace) {
      let workspace =
        server.schema.workspaces.findBy({ name: 'default' }) || server.create('workspace', 'default');
      statusReport.update('workspace', workspace);
    }
  },

  status: association('minutes-old', 'success'),

  unknown: trait({
    afterCreate(statusReport, server) {
      server.create('health', 'unknown', { statusReport });
    },
  }),

  alive: trait({
    afterCreate(statusReport, server) {
      server.create('health', 'alive', { statusReport });
    },
  }),

  ready: trait({
    afterCreate(statusReport, server) {
      server.create('health', 'ready', { statusReport });
    },
  }),

  down: trait({
    afterCreate(statusReport, server) {
      server.create('health', 'down', { statusReport });
    },
  }),

  partial: trait({
    afterCreate(statusReport, server) {
      server.create('health', 'partial', { statusReport });
    },
  }),

  'with-deployment-resources': trait({
    afterCreate(statusReport, server) {
      server.create('resource', 'random-deployment', { statusReport });
    },
  }),

  'with-release-resources': trait({
    afterCreate(statusReport, server) {
      server.create('resource', 'random-service', { statusReport });
    },
  }),
});
