/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { currentURL, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';

import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupSession } from 'waypoint/tests/helpers/login';

module('Acceptance | release detail', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);
  setupSession(hooks);

  test('displays a status report badge if available', async function (assert) {
    let project = this.server.create('project', { name: 'acme-project' });
    let application = this.server.create('application', { name: 'acme-app', project });
    let release = this.server.create('release', 'random', { application });
    this.server.create('status-report', 'ready', { application, target: release });

    await visit(`/default/acme-project/app/acme-app/release/seq/${release.sequence}`);

    assert.dom('[data-test-status-report-indicator="ready"]').exists();
  });

  test('displays no status report badge if none is available', async function (assert) {
    let project = this.server.create('project', { name: 'acme-project' });
    let application = this.server.create('application', { name: 'acme-app', project });
    let release = this.server.create('release', 'random', { application });

    await visit(`/default/acme-project/app/acme-app/release/seq/${release.sequence}`);

    assert.dom('[data-test-status-report-indicator]').doesNotExist();
  });

  test('redirects from the ID to the sequence route', async function (assert) {
    let project = this.server.create('project', { name: 'acme-project' });
    let application = this.server.create('application', { name: 'acme-app', project });
    let release = this.server.create('release', 'random', { application });

    await visit(`/default/acme-project/app/acme-app/release/${release.id}`);

    assert.equal(currentURL(), `/default/acme-project/app/acme-app/release/seq/${release.sequence}`);
  });

  test('displays the URL where appropriate', async function (assert) {
    let project = this.server.create('project', { name: 'acme-project' });
    let application = this.server.create('application', { name: 'acme-app', project });
    // Old superceded release
    this.server.create('release', 'random', {
      application,
      sequence: 1,
      state: 'CREATED',
      url: 'https://acme-app.test',
    });
    // New active release
    this.server.create('release', 'random', {
      application,
      sequence: 2,
      state: 'CREATED',
      url: 'https://acme-app.test',
    });

    await visit('/default/acme-project/app/acme-app/release/seq/1');

    assert.dom('[data-test-release-url]').doesNotExist('Doesn’t show the URL for an old release');

    await visit('/default/acme-project/app/acme-app/release/seq/2');

    assert
      .dom('[data-test-release-url]')
      .hasText('https://acme-app.test', 'Shows the URL for the active release');
  });
});
