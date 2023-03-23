/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { visit } from '@ember/test-helpers';
import { module, test } from 'qunit';

import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupSession } from 'waypoint/tests/helpers/login';

module('Acceptance | deployment detail', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);
  setupSession(hooks);

  test('displays a status report badge if available', async function (assert) {
    let project = this.server.create('project', { name: 'acme-project' });
    let application = this.server.create('application', { name: 'acme-app', project });
    let deployment = this.server.create('deployment', 'random', { application });
    this.server.create('status-report', 'ready', { application, target: deployment });

    await visit(`/default/acme-project/app/acme-app/deployments/${deployment.sequence}`);

    assert.dom('[data-test-status-report-indicator="ready"]').exists();
  });

  test('displays no status report badge if none is available', async function (assert) {
    let project = this.server.create('project', { name: 'acme-project' });
    let application = this.server.create('application', { name: 'acme-app', project });
    let deployment = this.server.create('deployment', 'random', { application });

    await visit(`/default/acme-project/app/acme-app/deployments/${deployment.sequence}`);

    assert.dom('[data-test-status-report-indicator]').doesNotExist();
  });
});
