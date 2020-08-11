import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import loadingMessage from 'cloud-ui-core/test-support/pages/components/loading/message';

module('Integration | Component | Loading/Message', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<Loading::Message>Description</Loading::Message>`);

    assert.dom(loadingMessage.containerSelector).exists();
    assert.dom(loadingMessage.containerSelector).containsText('Description');
    assert.dom(loadingMessage.containerSelector).hasClass('loadingState__message');
  });
});
