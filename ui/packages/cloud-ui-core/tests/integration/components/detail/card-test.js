import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { create } from 'ember-cli-page-object';
import cardPageObject from 'cloud-ui-core/test-support/pages/components/detail/card';

const card = create(cardPageObject);

module('Integration | Component | detail/card', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`
      <Detail::Card as |DC|>
        <DC.Header>Header Content</DC.Header>
        <DC.Content>Body Content</DC.Content>
      </Detail::Card>
    `);

    assert.ok(card.showsContainer, `the container renders`);
    assert.ok(card.showsHeader, `the header renders`);
    assert.equal(card.headerText, 'Header Content', `the header yields`);
    assert.ok(card.showsContent, `the content renders`);
    assert.equal(card.contentText, 'Body Content', `the content yields`);
  });
});
