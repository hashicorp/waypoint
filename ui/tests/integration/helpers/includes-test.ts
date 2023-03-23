/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Helper | includes', function (hooks) {
  setupRenderingTest(hooks);

  test('with a matching string', async function (assert) {
    await render(hbs`
      {{#if (includes "example" "exam")}}
        <div data-test-ok />
      {{/if}}
    `);

    assert.dom('[data-test-ok]').exists();
  });

  test('with a non-matching string', async function (assert) {
    await render(hbs`
      {{#if (includes "example" "test")}}
        <div data-test-ok />
      {{/if}}
    `);

    assert.dom('[data-test-ok]').doesNotExist();
  });

  test('with a matching array', async function (assert) {
    await render(hbs`
      {{#if (includes (array 1 2 3) 2)}}
        <div data-test-ok />
      {{/if}}
    `);

    assert.dom('[data-test-ok]').exists();
  });

  test('with a non-matching array', async function (assert) {
    await render(hbs`
      {{#if (includes (array 1 2 3) 4)}}
        <div data-test-ok />
      {{/if}}
    `);

    assert.dom('[data-test-ok]').doesNotExist();
  });

  test('with no haystack', async function (assert) {
    await render(hbs`
      {{#if (includes this.haystack "exam")}}
        <div data-test-ok />
      {{/if}}
    `);

    assert.dom('[data-test-ok]').doesNotExist();
  });

  test('with no needle', async function (assert) {
    await render(hbs`
      {{#if (includes "example")}}
        <div data-test-ok />
      {{/if}}
    `);

    assert.dom('[data-test-ok]').doesNotExist();
  });
});
