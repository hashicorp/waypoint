import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { create } from 'ember-cli-page-object';
import page from 'cloud-ui-core/test-support/pages/components/alert-banner';

const component = create(page);

module('Integration | Component | Alert Banner', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`
      <AlertBanner>
        <:title>Page Title</:title>
        <:content>Some content message</:content>
      </AlertBanner>
    `);

    assert.ok(component.isPresent, 'renders the alert banner');
    assert.ok(component.defaultStyle, 'renders default variant style');
  });
});
