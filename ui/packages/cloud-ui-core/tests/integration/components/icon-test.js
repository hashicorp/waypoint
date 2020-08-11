import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import waitForError from 'dummy/tests/helpers/wait-for-error';
import { SIZE_SCALE } from 'dummy/components/icon/consts';

module('Integration | Component | Icon', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<Icon />`);
    assert.dom('.icon').exists('renders root element');

    await render(hbs`<Icon @type='chevron-left' />`);
    assert.dom('.icon svg').exists('inlines the SVG');
  });

  test('it errors on unrecognized size', async function(assert) {
    let promise = waitForError();
    render(hbs`<Icon @size="no"/>`);
    let err = await promise;
    assert.ok(err.message.includes('@size for '), "errors when passed a size that's not allowed");
  });

  test('it adds size classes', async function(assert) {
    assert.expect(SIZE_SCALE.length);
    for (let size of SIZE_SCALE) {
      this.set('size', size);
      await render(hbs`<Icon @size={{size}} />`);
      assert.dom('.icon').hasClass(`icon--${size}`, `adds the size class for ${size}`);
    }
  });

  test('it forwards attributes', async function(assert) {
    await render(hbs`<Icon class="ah" aria-hidden="true" />`);
    assert.dom('.ah').hasAttribute('aria-hidden', 'true', 'renders class and aria-hidden');

    await render(hbs`<Icon class="al" aria-label="Testing" />`);
    assert.dom('.al').hasAttribute('aria-label', 'Testing', 'renders aria-label');
  });
});
