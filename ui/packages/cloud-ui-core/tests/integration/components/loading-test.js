import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import loading from 'cloud-ui-core/test-support/pages/components/loading';

module('Integration | Component | Loading', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`
      <Loading id="someId">Content</Loading>
    `);

    assert.dom(loading.containerSelector).exists();
    assert.dom(loading.containerSelector).containsText('Content');
    assert.dom(loading.containerSelector).hasClass('loadingState');
    assert.dom(loading.containerSelector).hasAttribute('id', 'someId');
    assert.dom(loading.iconContainerSelector).exists();
    assert.dom(loading.iconSelector).exists();
  });
});
