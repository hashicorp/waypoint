import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | status-badge', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders different states', async function(assert) {
    // Set any properties with this.set('myProperty', 'value');
    // Handle any actions with this.set('myAction', function(val) { ... });

    let errorBuild = {
      status: {
        state: 3,
      },
    };

    let successBuild = {
      status: {
        state: 2,
      },
    };

    let runningBuild = {
      status: {
        state: 1,
      },
    };

    let unknownBuild = {
      status: {
        state: 0,
      },
    };

    this.build = errorBuild;

    await render(hbs`<StatusBadge @model={{this.build}}/>`);

    assert.equal(this.element.getElementsByClassName('badge-status--error').length, 1);

    this.build = successBuild;

    await render(hbs`<StatusBadge @model={{this.build}}/>`);

    assert.equal(this.element.getElementsByClassName('badge-status--success').length, 1);

    this.build = runningBuild;

    await render(hbs`<StatusBadge @model={{this.build}}/>`);

    assert.equal(this.element.getElementsByClassName('badge-status--running').length, 1);

    this.build = unknownBuild;

    await render(hbs`<StatusBadge @model={{this.build}}/>`);

    assert.equal(this.element.getElementsByClassName('badge-status--unknown').length, 1);
  });
});
