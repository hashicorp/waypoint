import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import zeroStateHeader from 'cloud-ui-core/test-support/pages/components/zero-state/header';

module('Integration | Component | ZeroState/Header', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<ZeroState::Header>Header</ZeroState::Header>`);

    assert.dom(zeroStateHeader.containerSelector).exists();
    assert.dom(zeroStateHeader.containerSelector).containsText('Header');
    assert.dom(zeroStateHeader.containerSelector).hasClass('zero-state__header');
  });
});
