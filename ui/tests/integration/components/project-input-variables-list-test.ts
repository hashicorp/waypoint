import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn, click, pauseTest } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { create, collection, clickable, isPresent, fillable } from 'ember-cli-page-object';

const page = create({
  hasForm: isPresent('[data-test-input-variables-form]'),
  variablesList: collection('[data-test-input-variables-list-item]', {
    dropdown: clickable('[data-test-input-variables-dropdown]'),
    dropdownDelete: clickable('[data-test-input-variables-dropdown-delete]'),
  }),
  createButton: clickable('[data-test-input-variables-add-variable]'),
  cancelButton: clickable('[data-test-input-variables-edit-cancel]'),
  saveButton: clickable('[data-test-input-variables-edit-save]'),
  varName: fillable('[data-test-input-variables-var-name]'),
  varStr: fillable('[data-test-input-variables-var-str]'),
});

module('Integration | Component | project-input-variables-list', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    let project = {
      variablesList: [
        {
          name: 'Varname',
          str: 'foo',
        },
        {
          name: 'Varname2',
          str: 'foo2',
        },
      ],
    };
    this.set('project', project);
    await render(hbs`<ProjectInputVariables::List @project={{this.project}}/>`);
    assert.dom('.project-input-variables-list').exists('The list renders');
    assert.equal(page.variablesList.length, 2, 'the list contains all variables');
    await page.createButton();
    assert.ok(page.hasForm);
    await page.cancelButton();
    assert.equal(
      page.variablesList.length,
      2,
      'Attempt to create then cancel: the list still has the normal count of variables'
    );
    await page.createButton();
    await page.varName('var_name');
    await page.varStr('foozbarz');
    await page.saveButton();
    assert.equal(page.variablesList.length, 3, 'Create Variable: the list has the new variable');
    await page.variablesList.objectAt(0).dropdown();
    await page.variablesList.objectAt(0).dropdownDelete();
    assert.equal(page.variablesList.length, 2, 'Delete Variable: the variable has been removed');
  });
});
