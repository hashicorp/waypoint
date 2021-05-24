import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | project-input-variables-list', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    let project = {
      variablesList: [
        {
          name: 'Varname',
          str: 'foo',
        },
      ],
    };
    this.project = project;
    await render(hbs`<ProjectInputVariables::List @project={{this.project}}/>`);
    assert.equal(this.element.getElementsByClassName('project-input-variables-list').length, 1);
  });
});
