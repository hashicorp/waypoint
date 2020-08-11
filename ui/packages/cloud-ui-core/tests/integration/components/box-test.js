import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import waitForError from 'dummy/tests/helpers/wait-for-error';
import box from 'cloud-ui-core/test-support/pages/components/box';
import { DIMENSIONS_SIZE_SCALE, PADDING_SIZE_SCALE } from 'dummy/components/box/consts';

module('Integration | Component | box', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<Box />`);

    assert.dom(box.containerSelector).exists();
    assert.dom(box.containerSelector).hasClass('box--padding-top-sm');
    assert.dom(box.containerSelector).hasClass('box--padding-right-sm');
    assert.dom(box.containerSelector).hasClass('box--padding-bottom-sm');
    assert.dom(box.containerSelector).hasClass('box--padding-left-sm');
  });

  test('it renders: throws invalid size', async function(assert) {
    let promise = waitForError();
    render(hbs`<Box @padding="3xl" />`);
    let err = await promise;
    assert.ok(err.message.includes('@padding size for '), "errors when passed a size that's not allowed");
  });

  test('it renders: @padding xs', async function(assert) {
    await render(hbs`<Box @padding="xs" />`);

    assert.dom(box.containerSelector).exists();
    assert.dom(box.containerSelector).hasClass('box--padding-top-xs');
    assert.dom(box.containerSelector).hasClass('box--padding-right-xs');
    assert.dom(box.containerSelector).hasClass('box--padding-bottom-xs');
    assert.dom(box.containerSelector).hasClass('box--padding-left-xs');
  });

  test('it renders: two arguments', async function(assert) {
    await render(hbs`<Box @padding="xs md" />`);

    assert.dom(box.containerSelector).exists();
    assert.dom(box.containerSelector).hasClass('box--padding-top-xs');
    assert.dom(box.containerSelector).hasClass('box--padding-right-md');
    assert.dom(box.containerSelector).hasClass('box--padding-bottom-xs');
    assert.dom(box.containerSelector).hasClass('box--padding-left-md');
  });

  test('it renders: three arguments', async function(assert) {
    await render(hbs`<Box @padding="xs md lg" />`);

    assert.dom(box.containerSelector).exists();
    assert.dom(box.containerSelector).hasClass('box--padding-top-xs');
    assert.dom(box.containerSelector).hasClass('box--padding-right-md');
    assert.dom(box.containerSelector).hasClass('box--padding-bottom-lg');
    assert.dom(box.containerSelector).hasClass('box--padding-left-md');
  });

  test('it renders: four arguments', async function(assert) {
    await render(hbs`<Box @padding="xs md lg 2xl" />`);

    assert.dom(box.containerSelector).exists();
    assert.dom(box.containerSelector).hasClass('box--padding-top-xs');
    assert.dom(box.containerSelector).hasClass('box--padding-right-md');
    assert.dom(box.containerSelector).hasClass('box--padding-bottom-lg');
    assert.dom(box.containerSelector).hasClass('box--padding-left-2xl');
  });

  test('it renders: all @padding classes', async function(assert) {
    assert.expect(PADDING_SIZE_SCALE.length * DIMENSIONS_SIZE_SCALE.length * 4);
    for (let padding of PADDING_SIZE_SCALE) {
      this.set('padding', padding);
      await render(hbs`<Box @padding={{padding}} />`);
      for (let dimension of DIMENSIONS_SIZE_SCALE) {
        assert
          .dom(box.containerSelector)
          .hasClass(`box--padding-${dimension}-${padding}`, `adds the padding class for ${padding}`);
        assert
          .dom(box.containerSelector)
          .hasClass(`box--padding-${dimension}-${padding}`, `adds the padding class for ${padding}`);
        assert
          .dom(box.containerSelector)
          .hasClass(`box--padding-${dimension}-${padding}`, `adds the padding class for ${padding}`);
        assert
          .dom(box.containerSelector)
          .hasClass(`box--padding-${dimension}-${padding}`, `adds the padding class for ${padding}`);
      }
    }
  });
});
