import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | Typography', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`
      <Typography>
        template block text
      </Typography>
    `);

    assert.dom('p').exists();
    assert.dom('p').hasClass('body1');
  });

  test('it renders: @variant="h1"', async function(assert) {
    await render(hbs`
      <Typography @variant="h1">
        template block text
      </Typography>
    `);

    assert.dom('h1').exists();
  });

  test('it renders: @variant="h2"', async function(assert) {
    await render(hbs`
      <Typography @variant="h2">
        template block text
      </Typography>
    `);

    assert.dom('h2').exists();
  });

  test('it renders: @variant="h4"', async function(assert) {
    await render(hbs`
      <Typography @variant="h4">
        template block text
      </Typography>
    `);

    assert.dom('h4').exists();
  });

  test('it renders: @variant="h5"', async function(assert) {
    await render(hbs`
      <Typography @variant="h5">
        template block text
      </Typography>
    `);

    assert.dom('h5').exists();
  });

  test('it renders: @variant="h6"', async function(assert) {
    await render(hbs`
      <Typography @variant="h6">
        template block text
      </Typography>
    `);

    assert.dom('h6').exists();
  });

  test('it renders: @variant="subtitle1"', async function(assert) {
    await render(hbs`
      <Typography @variant="subtitle1">
        template block text
      </Typography>
    `);

    assert.dom('h6').exists();
    assert.dom('h6').hasClass('subtitle1');
  });

  test('it renders: @variant="subtitle2"', async function(assert) {
    await render(hbs`
      <Typography @variant="subtitle2">
        template block text
      </Typography>
    `);

    assert.dom('h6').exists();
    assert.dom('h6').hasClass('subtitle2');
  });

  test('it renders: @variant="body1"', async function(assert) {
    await render(hbs`
      <Typography @variant="body1">
        template block text
      </Typography>
    `);

    assert.dom('p').exists();
    assert.dom('p').hasClass('body1');
  });

  test('it renders: @variant="body2"', async function(assert) {
    await render(hbs`
      <Typography @variant="body2">
        template block text
      </Typography>
    `);

    assert.dom('p').exists();
    assert.dom('p').hasClass('body2');
  });

  test('it renders: @component="h1" @variant="h2"', async function(assert) {
    await render(hbs`
      <Typography @component="h1" @variant="h2">
        template block text
      </Typography>
    `);

    assert.dom('h1').exists();
    assert.dom('h1').hasClass('h2');
  });

  test('it renders: splattributes class', async function(assert) {
    await render(hbs`
      <Typography class="customClass">
        template block text
      </Typography>
    `);

    assert.dom('p').exists();
    assert.dom('p').hasClass('customClass');
  });
});
