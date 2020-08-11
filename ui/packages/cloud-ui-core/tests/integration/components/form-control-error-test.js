import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { create } from 'ember-cli-page-object';
import page from 'cloud-ui-core/test-support/pages/components/form-control-error';

const component = create(page);

module('Integration | Component | Form Control Error', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`
      <FormControlError>
        <:message>An error occurred</:message>
      </FormControlError>
    `);

    assert.ok(component.isPresent, 'renders the form control component');
  });
});
