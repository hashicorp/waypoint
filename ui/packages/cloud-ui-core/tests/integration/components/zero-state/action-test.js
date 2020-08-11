import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import zeroStateAction from 'cloud-ui-core/test-support/pages/components/zero-state/action';

module('Integration | Component | ZeroState/Action', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<ZeroState::Action>Action</ZeroState::Action>`);

    assert.dom(zeroStateAction.containerSelector).exists();
    assert.dom(zeroStateAction.containerSelector).containsText('Action');
    assert.dom(zeroStateAction.containerSelector).hasClass('zero-state__action');
  });
});
