/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';

import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupSession } from 'waypoint/tests/helpers/login';
import { visit } from '@ember/test-helpers';

module('Acceptance | deployment resource detail', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);
  setupSession(hooks);

  test('happy path', async function (assert) {
    let project = this.server.create('project', { name: 'my-project' });
    let application = this.server.create('application', { project, name: 'my-app' });
    let deployment = this.server.create('deployment', 'random', { application, sequence: 1 });
    let statusReport = this.server.create('status-report', 'ready', { application, target: deployment });
    let resource = this.server.create('resource', {
      statusReport,
      name: 'example-pod',
      state: {
        example: 'OK',
        image: 'example:1',
      },
    });

    await visit(
      `/default/${project.name}/app/${application.name}/deployments/${deployment.sequence}/resources/${resource.id}`
    );

    assert.dom('h1').containsText('example-pod');
    assert.dom('[data-test-resource-detail]').containsText('example 1');
    assert.dom('[data-test-json-viewer]').containsText('"example": "OK"');
  });

  test('error state', async function (assert) {
    let project = this.server.create('project', { name: 'my-project' });
    let application = this.server.create('application', { project, name: 'my-app' });
    let deployment = this.server.create('deployment', 'random', { application, sequence: 1 });

    await visit(
      `/default/${project.name}/app/${application.name}/deployments/${deployment.sequence}/resources/nope`
    );

    assert.dom('main').containsText('Resource nope not found');
  });
});
