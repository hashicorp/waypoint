import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import paper from 'cloud-ui-core/test-support/pages/components/paper';

module('Integration | Component | paper', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<Paper>content</Paper>`);

    assert.dom(paper.containerSelector).exists();
    assert.dom(paper.containerSelector).containsText('content');
  });

  test('it renders: @variant="oulined"', async function(assert) {
    await render(hbs`<Paper @variant="oulined">content</Paper>`);

    assert.dom(paper.containerSelector).exists();
    assert.dom(paper.containerSelector).containsText('content');
    assert.dom(paper.containerSelector).hasClass('oulined');
  });

  test('it renders: @square={{true}}', async function(assert) {
    await render(hbs`<Paper @square={{true}}>content</Paper>`);

    assert.dom(paper.containerSelector).exists();
    assert.dom(paper.containerSelector).containsText('content');
    assert.dom(paper.containerSelector).hasClass('square');
  });
});
