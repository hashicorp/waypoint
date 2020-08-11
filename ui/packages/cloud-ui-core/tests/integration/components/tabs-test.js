import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import tabs from 'cloud-ui-core/test-support/pages/components/tabs';

module('Integration | Component | tabs', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`
      <Tabs role="tablist" as |T|>
        <a href="#">Link</a>
        <a href="#">Link2</a>
      </Tabs>
    `);

    assert.dom(tabs.containerSelector).exists();
    assert.dom(tabs.containerSelector).hasAttribute('role', 'tablist');
  });
});
