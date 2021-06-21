import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { create, collection, clickable, isPresent, fillable } from 'ember-cli-page-object';

const page = create({
  hasForm: isPresent('[data-test-input-variables-form]'),
  variablesList: collection('[data-test-input-variables-list-item]', {
    dropdown: clickable('[data-test-input-variables-dropdown]'),
    dropdownDelete: clickable('[data-test-input-variables-dropdown-delete]'),
    isHcl: isPresent('[data-test-input-variables-list-item-is-hcl]'),
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
          hcl: 'hclval',
        },
      ],
    };
    this.set('project', project);
    await render(hbs`<ProjectInputVariables::List @project={{this.project}}/>`);
    assert.dom('.project-input-variables-list').exists('The list renders');
    assert.equal(page.variablesList.length, 2, 'the list contains all variables');
    assert.notOk(page.variablesList.objectAt(0).isHcl, 'the list contains a string variable');
    assert.ok(page.variablesList.objectAt(1).isHcl, 'the list contains a hcl variable');
  });

  test('the list can be edited and updated', async function (assert) {
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
    assert.ok(page.hasForm), 'Attempt to create: the form appears when the Add Variable button is clicked';
    await page.cancelButton();
    assert.equal(
      page.variablesList.length,
      2,
      'Attempt to create: the list still has the normal count of variables after cancelling'
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
