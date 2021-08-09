import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { create, collection, clickable, isPresent, fillable, text } from 'ember-cli-page-object';

const page = create({
  hasForm: isPresent('[data-test-config-variables-form]'),
  variablesList: collection('[data-test-config-variables-list-item]', {
    dropdown: clickable('[data-test-config-variables-dropdown]'),
    dropdownEdit: clickable('[data-test-config-variables-dropdown-edit]'),
    dropdownDelete: clickable('[data-test-config-variables-dropdown-delete]'),
    varName: text('[data-test-config-variables-var-name]'),
    varValue: text('[data-test-config-variables-var-value]'),
  }),
  createButton: clickable('[data-test-config-variables-add-variable]'),
  cancelButton: clickable('[data-test-config-variables-edit-cancel]'),
  saveButton: clickable('[data-test-config-variables-edit-save]'),
  varName: fillable('[data-test-config-variables-var-name]'),
  varStr: fillable('[data-test-config-variables-var-str]'),
});

module('Integration | Component | project-config-variables-list', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  test('it renders', async function (assert) {
    let dbproj = await this.server.create('project', { name: 'Proj1' });
    let dbVariablesList = this.server.createList('config-variable', 10, 'random');
    let varList = dbVariablesList.map(v => {
      return v.toProtobuf().toObject();
    });
    this.set('variablesList', varList);
    await render(hbs`<ProjectConfigVariables::List @variablesList={{this.variablesList}}/>`);
    assert.dom('.variables-list').exists('The list renders');
    assert.equal(page.variablesList.length, 10, 'it renders: the list has the proper length');
  });

  // test('only static variables are editable', async function (assert) {});
  // test('internal variables', async function (assert) {});
  // test('nameIsPath works', async function (assert) {});
});
