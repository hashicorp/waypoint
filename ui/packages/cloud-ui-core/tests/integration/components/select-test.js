import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { create } from 'ember-cli-page-object';
import pageObject from 'cloud-ui-core/test-support/pages/components/select';

const component = create(pageObject);

module('Integration | Component | select', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<Select />`);
    assert.ok(component.isRendered);
  });

  test('it renders options', async function(assert) {
    this.value = 'two';
    this.options = ['one', 'two', 'three'];
    await render(hbs`<Select @options={{this.options}} @value={{this.value}} />`);
    assert.equal(component.selectedOption().value, this.value);
    assert.equal(component.value, this.value);
  });

  test('it renders options from a list of objects', async function(assert) {
    this.value = 'second';
    this.options = [
      { number: '1', ordinal: 'first' },
      { number: '2', ordinal: 'second' },
      { number: '3', ordinal: 'third' },
    ];
    await render(hbs`<Select
      @options={{this.options}}
      @value={{this.value}}
      @valuePath='ordinal'
      @labelPath='number'
      />`);
    assert.equal(component.selectedOption().value, this.value);
    assert.equal(component.selectedOption().label, '2');
    assert.equal(component.value, this.value);
  });

  test('it sets value on change', async function(assert) {
    this.value = 'two';
    this.options = ['one', 'two', 'three'];
    this.setVal = (e) => {
      this.set('value', e.target.value);
    }
    await render(hbs`<Select
      @options={{this.options}}
      @value={{this.value}}
      {{on 'change' this.setVal}}
      />
      `);
    await component.fill('one');
    assert.equal(component.selectedOption().value, 'one');
    assert.equal(component.value, 'one');
  });

  test('it sets value on change with list of objects', async function(assert) {
    this.value = 'second';
    this.options = [
      { number: '1', ordinal: 'first' },
      { number: '2', ordinal: 'second' },
      { number: '3', ordinal: 'third' },
    ];
    this.setVal = (e) => {
      this.set('value', e.target.value);
    }
    await render(hbs`<Select
      @options={{this.options}}
      @value={{this.value}}
      @valuePath='ordinal'
      @labelPath='number'
      {{on 'change' this.setVal}}
      />`);
    await component.fill('first');
    assert.equal(component.selectedOption().value, 'first');
    assert.equal(component.selectedOption().label, '1');
    assert.equal(component.value, 'first');
  });
});
