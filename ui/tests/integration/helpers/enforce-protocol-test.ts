/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Helper | enforce-protocol', function (hooks) {
  setupRenderingTest(hooks);

  test('it adds a protocol to non-protocol urls', async function (assert) {
    this.set('inputValue', 'some-link.hacker.xyz');

    await render(hbs`{{enforce-protocol this.inputValue}}`);

    assert.equal(this.element.textContent?.trim(), 'https://some-link.hacker.xyz');
  });

  test('it keeps the protocol on http urls', async function (assert) {
    this.set('inputValue', 'http://some-link.hacker.xyz');

    await render(hbs`{{enforce-protocol this.inputValue}}`);

    assert.equal(this.element.textContent?.trim(), 'http://some-link.hacker.xyz');
  });

  test('it keeps the protocol on https urls', async function (assert) {
    this.set('inputValue', 'https://some-link.hacker.xyz');

    await render(hbs`{{enforce-protocol this.inputValue}}`);

    assert.equal(this.element.textContent?.trim(), 'https://some-link.hacker.xyz');
  });
});
