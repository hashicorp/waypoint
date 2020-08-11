import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { create } from 'ember-cli-page-object';
import page from 'cloud-ui-core/test-support/pages/components/toolbar';

let component = create(page);

module('Integration | Component | toolbar', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders with Actions', async function(assert) {
    await render(hbs`
      <Toolbar as |T|>
        <T.Actions>Ohai</T.Actions>
      </Toolbar>
    `);
    assert.ok(component.rendersToolbar, 'renders the toolbar');
    assert.ok(component.rendersActions, 'renders the toolbar');
    assert.notOk(component.rendersFilters, 'does not render the toolbar');
  });

  test('it renders with Filters', async function(assert) {
    await render(hbs`
      <Toolbar as |T|>
        <T.Filters>Ohai</T.Filters>
      </Toolbar>
    `);
    assert.ok(component.rendersToolbar, 'renders the toolbar');
    assert.notOk(component.rendersActions, 'does not render the actions');
    assert.ok(component.rendersFilters, 'renders the filters');
  });
});
