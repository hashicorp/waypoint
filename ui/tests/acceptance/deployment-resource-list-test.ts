/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';

import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupSession } from 'waypoint/tests/helpers/login';
import { visit } from '@ember/test-helpers';

module('Acceptance | deployment resource list', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);
  setupSession(hooks);

  test('happy path', async function (assert) {
    let project = this.server.create('project', { name: 'my-project' });
    let application = this.server.create('application', { project, name: 'my-app' });
    let deployment = this.server.create('deployment', 'random', { application, sequence: 1 });
    let statusReport = this.server.create('status-report', 'ready', { application, target: deployment });
    let resource = this.server.create('resource', { statusReport, name: 'example-pod' });

    await visit(`/default/${project.name}/app/${application.name}/deployments/${deployment.sequence}`);
    assert.dom('[data-test-resources-table]').containsText('example-pod');
    assert
      .dom(
        `[href="/default/${project.name}/app/${application.name}/deployments/${deployment.sequence}/resources/${resource.id}"]`
      )
      .exists();
  });

  test('empty state', async function (assert) {
    let project = this.server.create('project', { name: 'my-project' });
    let application = this.server.create('application', { project, name: 'my-app' });
    let deployment = this.server.create('deployment', 'random', { application, sequence: 1 });
    this.server.create('status-report', 'ready', { application, target: deployment });

    await visit(`/default/${project.name}/app/${application.name}/deployments/${deployment.sequence}`);

    assert.dom('[data-test-resources-table]').doesNotExist();
  });
});
