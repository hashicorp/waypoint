/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | json-viewer', function (hooks) {
  setupRenderingTest(hooks);

  test('with valid @json', async function (assert) {
    this.set('json', '{ "example": "OK" }');

    await render(hbs`
      <JsonViewer @json={{this.json}} />
    `);

    assert.dom('[data-test-json-viewer]').containsText('"example": "OK"');
  });

  test('with valid @label', async function (assert) {
    this.set('json', '{ "example": "OK" }');

    await render(hbs`
      <JsonViewer @json={{this.json}} @label="some-label" />
    `);

    assert.dom('[aria-label="some-label"]').exists();
  });

  test('with undefined @json', async function (assert) {
    await render(hbs`
      <JsonViewer />
    `);

    assert.dom('[data-test-json-viewer]').containsText('No source JSON provided');
  });

  test('with invalid @json', async function (assert) {
    this.set('json', '{');

    await render(hbs`
      <JsonViewer @json={{this.json}} />
    `);

    assert.dom('[data-test-error-message]').exists();
  });
});
