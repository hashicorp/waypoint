/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Server } from 'ember-cli-mirage';

export default function (server: Server): void {
  server.create('project', 'marketing-public');
  server.create('project', 'mutable-deployments');
  server.create('project', 'example-nodejs');
}
