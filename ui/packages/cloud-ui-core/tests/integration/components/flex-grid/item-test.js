import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import flexGridItem from 'cloud-ui-core/test-support/pages/components/flex-grid/item';

module('Integration | Component | GridItem', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<FlexGrid::Item>content</FlexGrid::Item>`);

    assert.dom(flexGridItem.containerSelector).exists();
    assert.dom(flexGridItem.containerSelector).containsText('content');
  });

  test('it renders: @xs, @sm, @md, @lg props', async function(assert) {
    await render(hbs`
      <FlexGrid::Item @xs={{3}} @sm={{6}} @md={{3}} @lg={{6}}>
        content
      </FlexGrid::Item>
    `);

    assert.dom(flexGridItem.containerSelector).exists();
    assert.dom(flexGridItem.containerSelector).hasClass('col-xs-3');
    assert.dom(flexGridItem.containerSelector).hasClass('col-sm-6');
    assert.dom(flexGridItem.containerSelector).hasClass('col-md-3');
    assert.dom(flexGridItem.containerSelector).hasClass('col-lg-6');
  });

  test('it renders: @xsOffset, @smOffset, @mdOffset, @lgOffset props', async function(assert) {
    await render(hbs`
      <FlexGrid::Item @xsOffset={{3}} @smOffset={{6}} @mdOffset={{3}} @lgOffset={{6}}>
        content
      </FlexGrid::Item>
    `);

    assert.dom(flexGridItem.containerSelector).exists();
    assert.dom(flexGridItem.containerSelector).hasClass('col-xs-offset-3');
    assert.dom(flexGridItem.containerSelector).hasClass('col-sm-offset-6');
    assert.dom(flexGridItem.containerSelector).hasClass('col-md-offset-3');
    assert.dom(flexGridItem.containerSelector).hasClass('col-lg-offset-6');
  });
});
