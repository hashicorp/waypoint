import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import flexGrid from 'cloud-ui-core/test-support/pages/components/flex-grid';

module('Integration | Component | FlexGrid', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<FlexGrid>content</FlexGrid>`);

    assert.dom(flexGrid.containerSelector).exists();
    assert.dom(flexGrid.containerSelector).hasClass('row');
    assert.dom(flexGrid.containerSelector).containsText('content');
  });

  test('it renders: @reverse', async function(assert) {
    await render(hbs`<FlexGrid @reverse={{true}}>content</FlexGrid>`);

    assert.dom(flexGrid.containerSelector).exists();
    assert.dom(flexGrid.containerSelector).hasClass('row');
    assert.dom(flexGrid.containerSelector).hasClass('reverse');
    assert.dom(flexGrid.containerSelector).containsText('content');
  });
});
