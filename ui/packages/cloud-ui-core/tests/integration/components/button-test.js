import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import button from 'cloud-ui-core/test-support/pages/components/button';
import { DEFAULT_VARIANT, DEFAULT_VARIANT_MAPPING, VARIANT_SCALE } from 'dummy/components/button/consts';

module('Integration | Component | Button', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<Button>Create Network</Button>`);
    assert.dom(button.containerSelector).exists();
    assert.dom(button.containerSelector).hasClass('button');
    assert.dom(button.containerSelector).hasClass(DEFAULT_VARIANT_MAPPING[DEFAULT_VARIANT]);
  });

  test('it adds variant classes', async function(assert) {
    assert.expect(VARIANT_SCALE.length);
    for (let variant of VARIANT_SCALE) {
      this.set('variant', variant);
      await render(hbs`<Button @variant={{variant}} />`);
      assert
        .dom(button.containerSelector)
        .hasClass(DEFAULT_VARIANT_MAPPING[variant], `adds the variant class for ${variant}`);
    }
  });

  test('it adds compact class', async function(assert) {
    this.set('compact', true);
    await render(hbs`<Button @compact={{compact}} />`);
    assert.dom(button.containerSelector).hasClass('button--compact');
  });

  test('it forwards attributes', async function(assert) {
    await render(hbs`<Button type="submit" aria-hidden="true">Create Network</Button>`);
    assert.dom(button.containerSelector).hasAttribute('type', 'submit', 'renders the type');
    assert.dom(button.containerSelector).hasAttribute('aria-hidden', 'true', 'renders the aria-hidden');
  });
});
