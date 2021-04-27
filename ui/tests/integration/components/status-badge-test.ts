import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | status-badge', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders different states', async function(assert) {

    let partialBuild = {
      status: {
        state: 4,
      },
    };

    let downBuild = {
      status: {
        state: 3,
      },
    };

    let readyBuild = {
      status: {
        state: 2,
      },
    };

    let aliveBuild = {
      status: {
        state: 1,
      },
    };

    let unknownBuild = {
      status: {
        state: 0,
      },
    };

    this.build = partialBuild;

    await render(hbs`<StatusBadge @model={{this.build}}/>`);

    assert.equal(this.element.getElementsByClassName('badge-status--partial').length, 1);

    this.build = downBuild;

    await render(hbs`<StatusBadge @model={{this.build}}/>`);

    assert.equal(this.element.getElementsByClassName('badge-status--down').length, 1);

    this.build = readyBuild;

    await render(hbs`<StatusBadge @model={{this.build}}/>`);

    assert.equal(this.element.getElementsByClassName('badge-status--ready').length, 1);

    this.build = aliveBuild;

    await render(hbs`<StatusBadge @model={{this.build}}/>`);

    assert.equal(this.element.getElementsByClassName('badge-status--alive').length, 1);

    this.build = unknownBuild;

    await render(hbs`<StatusBadge @model={{this.build}}/>`);

    assert.equal(this.element.getElementsByClassName('badge-status--unknown').length, 1);
  });
});
