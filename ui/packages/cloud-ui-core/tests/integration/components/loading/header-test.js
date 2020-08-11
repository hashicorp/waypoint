import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import loadingHeader from 'cloud-ui-core/test-support/pages/components/loading/header';

module('Integration | Component | Loading/Header', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<Loading::Header>Header</Loading::Header>`);

    assert.dom(loadingHeader.containerSelector).exists();
    assert.dom(loadingHeader.containerSelector).containsText('Header');
    assert.dom(loadingHeader.containerSelector).hasClass('loadingState__header');
  });
});
