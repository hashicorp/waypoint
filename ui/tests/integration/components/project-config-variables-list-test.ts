/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { clickable, collection, create, fillable, isPresent, text } from 'ember-cli-page-object';
import { module, test } from 'qunit';
import { render } from '@ember/test-helpers';

import { TestContext } from 'ember-test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupRenderingTest } from 'ember-qunit';

const page = create({
  hasForm: isPresent('[data-test-config-variables-form]'),
  variablesList: collection('[data-test-config-variables-list-item]', {
    dropdown: clickable('[data-test-config-variables-dropdown]'),
    dropdownEdit: clickable('[data-test-config-variables-dropdown-edit]'),
    hasDropDown: isPresent('[data-test-config-variables-dropdown]'),
    hasDropDownEdit: isPresent('[data-test-config-variables-dropdown-edit]'),
    dropdownDelete: clickable('[data-test-config-variables-dropdown-delete]'),
    varName: text('[data-test-config-variables-var-name]'),
    varValue: text('[data-test-config-variables-var-value]'),
    varNameIsPath: text('[data-test-config-variables-var-name-is-path]'),
    varInternal: text('[data-test-config-variables-var-internal]'),
  }),
  createButton: clickable('[data-test-config-variables-add-variable]'),
  cancelButton: clickable('[data-test-config-variables-edit-cancel]'),
  saveButton: clickable('[data-test-config-variables-edit-save]'),
  varName: fillable('[data-test-config-variables-var-name]'),
  varStatic: fillable('[data-test-config-variables-var-static]'),
  varNameIsPath: clickable('[data-test-config-variables-name-is-path-toggle]'),
  varInternal: clickable('[data-test-config-variables-internal-toggle]'),
});

module('Integration | Component | project-config-variables-list', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function (this: TestContext) {
    // We have to register any types we expect to use in this component
    this.owner.lookup('service:flash-messages').registerTypes(['success', 'error']);
  });

  test('it renders', async function (assert) {
    let dbproj = await this.server.create('project', { name: 'Proj1' });
    let proj = dbproj.toProtobuf().toObject();
    let dbVariablesList = this.server.createList('config-variable', 10, 'random');
    let varList = dbVariablesList.map((v) => {
      return v.toProtobuf().toObject();
    });
    this.set('variablesList', varList);
    this.set('project', proj);
    await render(
      hbs`<ProjectConfigVariables::List @variablesList={{this.variablesList}} @project={{this.project}}/>`
    );
    assert.dom('.variables-list').exists('The list renders');
    assert.equal(page.variablesList.length, 10, 'it renders: the list has the proper length');
  });

  test('adding and deleting variables works', async function (assert) {
    let dbproj = await this.server.create('project', { name: 'Proj1' });
    let dbVariablesList = this.server.createList('config-variable', 3, 'random', { project: dbproj });
    let proj = dbproj.toProtobuf().toObject();
    let varList = dbVariablesList.map((v) => {
      return v.toProtobuf().toObject();
    });
    this.set('variablesList', varList);
    this.set('project', proj);
    await render(
      hbs`<ProjectConfigVariables::List @variablesList={{this.variablesList}} @project={{this.project}}/>`
    );

    assert.dom('.variables-list').exists('The list renders');
    assert.equal(page.variablesList.length, 3, 'the list contains all variables');
    await page.createButton();
    assert.ok(page.hasForm, 'Attempt to create: the form appears when the Add Variable button is clicked');
    await page.cancelButton();
    assert.notOk(page.hasForm, 'Attempt to create: the form is hidden after canceling');
    assert.equal(
      page.variablesList.length,
      3,
      'Attempt to create: the list still has the normal count of variables after cancelling'
    );
    await page.createButton();
    await page.varName('var_name');
    await page.varStatic('foozbarz');
    await page.varNameIsPath();
    await page.varInternal();
    await page.saveButton();
    assert.notOk(page.hasForm, 'Create Variable: the form disappears after creation');
    assert.equal(page.variablesList.length, 4, 'Create Variable: the list has the new variable');
    assert.equal(page.variablesList.objectAt(3).varName, 'var_name', 'Var name is correct');
    assert.equal(page.variablesList.objectAt(3).varValue, 'foozbarz', 'Var value is correct');
    assert.equal(page.variablesList.objectAt(3).varNameIsPath, 'true', 'name is path is correct');
    assert.equal(page.variablesList.objectAt(3).varInternal, 'true', 'internal is set correctly is correct');
  });

  test('only static variables are editable', async function (assert) {
    let dbproj = await this.server.create('project', { name: 'Proj1' });
    let proj = dbproj.toProtobuf().toObject();
    let dbVariablesList = this.server.createList('config-variable', 3, 'random', { project: dbproj });
    let dynamicVar = this.server.create('config-variable', 'dynamic', { project: dbproj });
    dbVariablesList.push(dynamicVar);
    let varList = dbVariablesList.map((v) => {
      return v.toProtobuf().toObject();
    });
    this.set('variablesList', varList);
    this.set('project', proj);
    await render(
      hbs`<ProjectConfigVariables::List @variablesList={{this.variablesList}} @project={{this.project}}/>`
    );
    assert.dom('.variables-list').exists('The list renders');

    assert.equal(page.variablesList.length, 4, 'the list contains all variables');
    await page.variablesList.objectAt(0).dropdown();
    assert.ok(page.variablesList.objectAt(0).hasDropDownEdit, 'Static Variable is editable');
    assert.notOk(page.variablesList.objectAt(3).hasDropDown, 'Dynamic Variable is not editable or deletable');
    await page.variablesList.objectAt(0).dropdownEdit();
    await page.varStatic('foozbarz');
    await page.saveButton();
    assert.notOk(page.hasForm, 'Create Variable: the form disappears after creation');
  });

  test('renaming variables works', async function (assert) {
    let dbproj = await this.server.create('project', { name: 'Proj1' });
    let proj = dbproj.toProtobuf().toObject();
    let dbVariablesList = this.server.createList('config-variable', 4, 'random', { project: dbproj });
    let varList = dbVariablesList.map((v) => {
      return v.toProtobuf().toObject();
    });
    this.set('variablesList', varList);
    this.set('project', proj);
    await render(
      hbs`<ProjectConfigVariables::List @variablesList={{this.variablesList}} @project={{this.project}}/>`
    );
    assert.equal(page.variablesList.length, 4, 'the list contains the right number of variables');
    await page.variablesList.objectAt(0).dropdown();
    await page.variablesList.objectAt(0).dropdownEdit();
    await page.varName('edited_var_name');
    await page.saveButton();
    assert.equal(
      page.variablesList.length,
      4,
      'the list contains the right number of variables after name edition'
    );
  });
});
