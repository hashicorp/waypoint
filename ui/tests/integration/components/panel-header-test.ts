/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | panel-header', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders an empty component when params not passed', async function (assert) {
    await render(hbs`<PanelHeader @artifact="build" />`);

    assert.dom('[data-test-panel-header]').containsText('Build');
  });

  test('it renders an empty component when params passed', async function (assert) {
    await render(hbs`
      <PanelHeader @artifact="build" @sequence={{3}}/>
    `);

    assert.dom('[data-test-panel-header]').exists();
    assert.dom('[data-test-panel-header]').containsText('Build');
    assert.dom('[data-test-panel-header]').containsText('v3');
  });

  test('it does not renders empty badge when sequence param missing', async function (assert) {
    await render(hbs`
      <PanelHeader @artifact="build"/>
    `);

    assert.dom('[data-test-panel-header]').exists();
    assert.dom('[data-test-panel-header]').containsText('Build');
    assert.dom('[data-test-panel-header]').doesNotContainText('v');
  });
});
