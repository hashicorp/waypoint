/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | section', function (hooks) {
  setupRenderingTest(hooks);

  test('basic rendering', async function (assert) {
    await render(hbs`
      <Section>
        <:heading>Heading</:heading>
        <:body>Body</:body>
      </Section>
    `);

    assert.dom('section').hasClass('section--expanded');
    assert.dom('section').containsText('Heading');
    assert.dom('section').containsText('Body');
  });

  test('toggling', async function (assert) {
    await render(hbs`
      <Section>
        <:heading>Heading</:heading>
        <:body>Body</:body>
      </Section>
    `);

    await click('[data-test-section-toggle]');

    assert.dom('section').doesNotContainText('Body');

    await click('[data-test-section-toggle]');

    assert.dom('section').containsText('Body');
  });

  test('with @expanded={{false}}', async function (assert) {
    await render(hbs`
      <Section @expanded={{false}}>
        <:heading>Heading</:heading>
        <:body>Body</:body>
      </Section>
    `);

    assert.dom('section').doesNotContainText('Body');
  });
});
