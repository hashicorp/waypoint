import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import zeroState from 'cloud-ui-core/test-support/pages/components/zero-state';

module('Integration | Component | ZeroState', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`
      <ZeroState>Content</ZeroState>
    `);

    assert.dom(zeroState.containerSelector).exists();
    assert.dom(zeroState.containerSelector).hasClass('zero-state');
    assert.dom(zeroState.containerSelector).containsText('Content');
  });
});
