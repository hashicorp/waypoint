import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { CODE_BLOCK_CONTAINER_SELECTOR } from 'cloud-ui-core/test-support/pages/components/code-block';

module('Integration | Component | code-block', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`
      <CodeBlock>
        some code
      </CodeBlock>
    `);

    assert.dom(CODE_BLOCK_CONTAINER_SELECTOR).hasClass('code-block');
    assert.dom(CODE_BLOCK_CONTAINER_SELECTOR).containsText('some code');
  });
});
