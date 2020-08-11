import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import {
  CONTAINER_SELECTOR,
  HEADER_SELECTOR,
  TITLE_SELECTOR,
  ZERO_STATE_SELECTOR,
} from 'cloud-ui-core/test-support/pages/components/detail/section';

module('Integration | Component | detail/section', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`
      <Detail::Section @title='Title'>
        Content
      </Detail::Section>
    `);

    assert.dom(CONTAINER_SELECTOR).exists();
    assert.dom(CONTAINER_SELECTOR).containsText('Content');
    assert.dom(HEADER_SELECTOR).exists();
    assert.dom(TITLE_SELECTOR).exists();
    assert.dom(TITLE_SELECTOR).containsText('Title');
  });

  test('it renders: with zeroState', async function(assert) {
    await render(hbs`
      <Detail::Section @title='Title' as |DS|>
        <DS.ZeroState>Some empty message</DS.ZeroState>
      </Detail::Section>
    `);

    assert.dom(ZERO_STATE_SELECTOR).exists();
    assert.dom(ZERO_STATE_SELECTOR).containsText('Some empty message');
  });
});
