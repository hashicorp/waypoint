/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import { TestContext } from 'ember-test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { create, collection, clickable, isPresent, fillable, text } from 'ember-cli-page-object';

const page = create({
  hasForm: isPresent('[data-test-input-variables-form]'),
  variablesList: collection('[data-test-input-variables-list-item]', {
    dropdown: clickable('[data-test-input-variables-dropdown]'),
    dropdownEdit: clickable('[data-test-input-variables-dropdown-edit]'),
    dropdownDelete: clickable('[data-test-input-variables-dropdown-delete]'),
    isHcl: isPresent('[data-test-input-variables-list-item-is-hcl]'),
    varName: text('[data-test-input-variables-var-name]'),
    varValue: text('[data-test-input-variables-var-value]'),
  }),
  createButton: clickable('[data-test-input-variables-add-variable]'),
  cancelButton: clickable('[data-test-input-variables-edit-cancel]'),
  saveButton: clickable('[data-test-input-variables-edit-save]'),
  varName: fillable('[data-test-input-variables-var-name]'),
  varStr: fillable('[data-test-input-variables-var-str]'),
});

module('Integration | Component | project-input-variables-list', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function (this: TestContext) {
    // We have to register any types we expect to use in this component
    this.owner.lookup('service:flash-messages').registerTypes(['success', 'error']);
  });

  test('it renders', async function (assert) {
    let dbproj = await this.server.create('project', 'with-input-variables', { name: 'Proj1' });
    let project = dbproj.toProtobuf();
    this.set('project', project.toObject());

    await render(hbs`<ProjectInputVariables::List @project={{this.project}}/>`);
    assert.dom('.variables-list').exists('The list renders');
    assert.equal(page.variablesList.length, 4, 'the list contains all variables');
    assert.notOk(page.variablesList.objectAt(0).isHcl, 'the list contains a string variable');
    assert.ok(page.variablesList.objectAt(2).isHcl, 'the list contains a hcl variable');
  });

  test('adding and deleting variables works', async function (assert) {
    let dbproj = await this.server.create('project', 'with-input-variables', { name: 'Proj3' });
    let project = dbproj.toProtobuf();
    this.set('project', project.toObject());

    await render(hbs`<ProjectInputVariables::List @project={{this.project}}/>`);
    assert.dom('.variables-list').exists('The list renders');
    assert.equal(page.variablesList.length, 4, 'the list contains all variables');
    await page.createButton();
    assert.ok(page.hasForm, 'Attempt to create: the form appears when the Add Variable button is clicked');
    await page.cancelButton();
    assert.equal(
      page.variablesList.length,
      4,
      'Attempt to create: the list still has the normal count of variables after cancelling'
    );
    await page.createButton();
    await page.varName('var_name');
    await page.varStr('foozbarz');
    await page.saveButton();
    assert.equal(page.variablesList.length, 5, 'Create Variable: the list has the new variable');
    await page.variablesList.objectAt(0).dropdown();
    await page.variablesList.objectAt(0).dropdownDelete();
    assert.equal(page.variablesList.length, 4, 'Delete Variable: the variable has been removed');
  });

  test('editing variables works', async function (assert) {
    let dbproj = await this.server.create('project', { name: 'Proj3' });
    let project = dbproj.toProtobuf();
    this.set('project', project.toObject());

    await render(hbs`<ProjectInputVariables::List @project={{this.project}}/>`);
    assert.dom('.variables-list').doesNotExist('the list is empty initially');
    assert.equal(page.variablesList.length, 0, 'the list contains no variables');
    await page.createButton();
    assert.ok(page.hasForm), 'Attempt to create: the form appears when the Add Variable button is clicked';
    await page.varName('var_name');
    await page.varStr('foozbarz');
    await page.saveButton();
    assert.equal(page.variablesList.length, 1, 'Create Variable: the list has the new variable');
    assert.equal(page.variablesList.objectAt(0).varName, 'var_name', 'the variable name is correct');
    assert.equal(page.variablesList.objectAt(0).varValue, 'foozbarz', 'the variable value is correct');
    await page.variablesList.objectAt(0).dropdown();
    await page.variablesList.objectAt(0).dropdownEdit();
    await page.varName('var_name_edited');
    await page.varStr('foozbarz_edited');
    await page.saveButton();
    assert.equal(page.variablesList.length, 1, 'The list has the edited variable');
    assert.equal(
      page.variablesList.objectAt(0).varName,
      'var_name_edited',
      'the updated variable name is correct'
    );
    assert.equal(
      page.variablesList.objectAt(0).varValue,
      'foozbarz_edited',
      'the updated variable value is correct'
    );
  });

  test('sensitive variables are hidden in list', async function (assert) {
    let dbproj = await this.server.create('project', { name: 'Proj3' });
    this.server.create('variable', 'is-sensitive', { project: dbproj });
    let project = dbproj.toProtobuf();
    this.set('project', project.toObject());

    await render(hbs`<ProjectInputVariables::List @project={{this.project}}/>`);

    assert.dom('[data-test-sensitive-var-badge]').exists();
  });

  test('sensitive variables are hidden in forms', async function (assert) {
    let dbproj = await this.server.create('project', { name: 'Proj3' });
    this.server.create('variable', 'is-sensitive', { project: dbproj });
    let project = dbproj.toProtobuf();
    this.set('project', project.toObject());

    await render(hbs`<ProjectInputVariables::List @project={{this.project}}/>`);

    await page.variablesList.objectAt(0).dropdown();
    await page.variablesList.objectAt(0).dropdownEdit();

    assert.dom('[data-test-input-variables-var-str]').hasAttribute('placeholder', 'sensitive - write only');
  });
});
