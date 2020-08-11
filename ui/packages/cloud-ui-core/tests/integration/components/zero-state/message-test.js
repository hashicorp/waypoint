import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import zeroStateMessage from 'cloud-ui-core/test-support/pages/components/zero-state/message';

module('Integration | Component | ZeroState/Message', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<ZeroState::Message>Description</ZeroState::Message>`);

    assert.dom(zeroStateMessage.containerSelector).exists();
    assert.dom(zeroStateMessage.containerSelector).containsText('Description');
    assert.dom(zeroStateMessage.containerSelector).hasClass('zero-state__message');
  });
});
