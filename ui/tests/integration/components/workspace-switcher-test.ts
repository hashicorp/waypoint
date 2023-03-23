/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { render, waitFor, settled } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupRenderingTest } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import { module, test } from 'qunit';
import { Response } from 'miragejs';

module('Integration | Component | workspace-switcher', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  test('happy path', async function (assert) {
    this.server.create('workspace', { name: 'default' });
    this.server.create('workspace', { name: 'production' });

    await render(hbs`
      <WorkspaceSwitcher
        @current="default"
        @route="workspace.projects.project.app.builds"
        @models={{array "test-project" "test-app"}}
        @isOpen={{true}}
      />
    `);

    assert.dom('[data-test-dropdown-trigger]').hasText('default');
    assert.dom('a[href="/default/test-project/app/test-app/builds"]').hasAttribute('aria-current', 'page');
    assert.dom('a[href="/production/test-project/app/test-app/builds"]').hasNoAttribute('aria-current');
  });

  test('while loading', async function (assert) {
    this.server.create('workspace', { name: 'default' });
    this.server.create('workspace', { name: 'production' });

    render(hbs`
      <WorkspaceSwitcher
        @current="default"
        @route="workspace.projects.project.app.builds"
        @models={{array "test-project" "test-app"}}
        @isOpen={{true}}
      />
      <div data-test-sentinel />
    `);

    await waitFor('[data-test-sentinel]');

    assert.dom('[data-test-workspace-switcher]').doesNotExist();

    await settled();

    assert.dom('[data-test-workspace-switcher]').exists();
  });

  test('when we fail to fetch workspaces from the API', async function (assert) {
    this.server.post('/ListWorkspaces', () => new Response(500));

    await render(hbs`
      <WorkspaceSwitcher
        @current="default"
        @route="workspace.projects.project.app.builds"
        @models={{array "test-project" "test-app"}}
        @isOpen={{true}}
      />
    `);

    assert.dom('[data-test-workspace-switcher-error-state]').exists();
  });

  test('with no workspaces', async function (assert) {
    await render(hbs`
      <WorkspaceSwitcher
        @current="default"
        @route="workspace.projects.project.app.builds"
        @models={{array "test-project" "test-app"}}
        @isOpen={{true}}
      />
    `);

    assert.dom('[data-test-workspace-switcher]').doesNotExist();
  });

  test('with only one workspace', async function (assert) {
    this.server.create('workspace', { name: 'default' });

    await render(hbs`
      <WorkspaceSwitcher
        @current="default"
        @route="workspace.projects.project.app.builds"
        @models={{array "test-project" "test-app"}}
        @isOpen={{true}}
      />
    `);

    assert.dom('[data-test-workspace-switcher]').doesNotExist();
  });

  test('with only one workspace, that differs from @current', async function (assert) {
    this.server.create('workspace', { name: 'production' });

    await render(hbs`
      <WorkspaceSwitcher
        @current="default"
        @route="workspace.projects.project.app.builds"
        @models={{array "test-project" "test-app"}}
        @isOpen={{true}}
      />
    `);

    assert.dom('[data-test-workspace-switcher]').exists();
  });
});
