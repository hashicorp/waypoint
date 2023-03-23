/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Helper | clean-up-url', function (hooks) {
  setupRenderingTest(hooks);

  test('it removes https:// from input', async function (assert) {
    this.set('inputValue', 'https://wildly-intent-honeybee.waypoint.run');

    await render(hbs`{{clean-up-url this.inputValue}}`);

    assert.equal(this.element.textContent?.trim(), 'wildly-intent-honeybee.waypoint.run');
  });

  test('it removes http:// from input', async function (assert) {
    this.set('inputValue', 'http://wildly-intent-honeybee.waypoint.run');

    await render(hbs`{{clean-up-url this.inputValue}}`);

    assert.equal(this.element.textContent?.trim(), 'wildly-intent-honeybee.waypoint.run');
  });
});
