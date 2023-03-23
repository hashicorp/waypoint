/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { click, visit, currentURL } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupApplicationTest } from 'ember-qunit';
import { module, test } from 'qunit';
import { setupSession } from '../helpers/login';
import { currentSession } from 'ember-simple-auth/test-support';

module('Acceptance | workspaces', function (hooks) {
  setupApplicationTest(hooks);
  setupSession(hooks);
  setupMirage(hooks);

  test('switching between workspaces while looking at deployments', async function (assert) {
    let staging = this.server.create('workspace', { name: 'staging' });
    let production = this.server.create('workspace', { name: 'production' });
    let project = this.server.create('project', { name: 'test-project' });
    let application = this.server.create('application', { name: 'test-app', project });

    this.server.create('deployment', 'random', { application, seq: 1, workspace: staging });
    this.server.create('deployment', 'random', { application, seq: 2, workspace: production });

    await visit(`/staging/test-project/app/test-app`);

    assert.dom('[data-test-workspace-switcher]').containsText('staging');

    await click('[data-test-workspace-switcher] [data-test-dropdown-trigger]');
    await click(`[data-test-workspace-link="production"]`);

    assert.equal(currentURL(), `/production/test-project/app/test-app/deployments/2`);
    assert.dom('[data-test-workspace-switcher]').containsText('production');
  });

  test('switching between workspaces while looking at builds', async function (assert) {
    let staging = this.server.create('workspace', { name: 'staging' });
    let production = this.server.create('workspace', { name: 'production' });
    let project = this.server.create('project', { name: 'test-project' });
    let application = this.server.create('application', { name: 'test-app', project });

    this.server.create('build', 'random', { application, seq: 1, workspace: staging });
    this.server.create('build', 'random', { application, seq: 2, workspace: production });

    await visit(`/staging/test-project/app/test-app/builds`);

    assert.dom('[data-test-workspace-switcher]').containsText('staging');

    await click('[data-test-workspace-switcher] [data-test-dropdown-trigger]');
    await click(`[data-test-workspace-link="production"]`);

    assert.equal(currentURL(), `/production/test-project/app/test-app/builds`);
    assert.dom('[data-test-workspace-switcher]').containsText('production');
  });

  test('selects workspace in local storage if valid', async function (assert) {
    this.server.create('workspace', { name: 'dev' });
    this.server.create('workspace', { name: 'production' });

    let session = currentSession();

    session.set('data.workspace', 'production');

    await visit('/');

    assert.equal(currentURL(), '/production');
  });

  test('selects default workspace if it exists', async function (assert) {
    this.server.create('workspace', { name: 'alpha' });
    this.server.create('workspace', { name: 'default' });

    await visit('/');

    assert.equal(currentURL(), '/default');
  });

  test('selects alphabetically first workspace if default does not exist', async function (assert) {
    this.server.create('workspace', { name: 'beta' });
    this.server.create('workspace', { name: 'alpha' });

    await visit('/');

    assert.equal(currentURL(), '/alpha');
  });

  test('selects default workspace if no concrete workspaces exist', async function (assert) {
    await visit('/');

    assert.equal(currentURL(), '/default');
  });

  test('selects the alphabetically first workspace if workspace in local storage is invalid', async function (assert) {
    this.server.create('workspace', { name: 'alpha' });

    let session = currentSession();

    session.set('data.workspace', 'nope');

    await visit('/');

    assert.equal(currentURL(), '/alpha');
  });

  test('remembers the current workspace', async function (assert) {
    this.server.create('workspace', { name: 'alpha' });
    this.server.create('workspace', { name: 'default' });

    await visit('/alpha');
    await visit('/');

    assert.equal(currentURL(), '/alpha');
  });

  test('forgets the workspace on logout', async function (assert) {
    this.server.create('workspace', { name: 'alpha' });
    this.server.create('workspace', { name: 'default' });

    let session = currentSession();

    await visit('/alpha');
    await click('[data-test-logout-button]');

    assert.equal(session.data.workspace, undefined, 'workspace no longer in session store');
  });
});
