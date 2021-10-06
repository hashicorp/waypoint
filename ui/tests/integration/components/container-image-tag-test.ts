import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | container-image-tag', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    this.set('statusReport', {
      resourcesList: [
        {
          type: 'container',
          stateJson: '{"Config": {"Image": "docker:tag"}}',
        },
      ],
    });

    await render(hbs`<ContainerImageTag @statusReport={{this.statusReport}}/>`);
    assert.equal(this.element?.textContent?.trim().replace(/\s+/, ' '), 'docker tag');
  });

  test('it renders n/a when a status report does not exist', async function (assert) {
    this.set('statusReport2', {});

    await render(hbs`<ContainerImageTag @statusReport={{this.statusReport2}}/>`);
    assert.equal(this.element?.textContent?.trim().replace(/\s+/, ' '), 'n/a');
  });

  test('it renders multiple tags', async function (assert) {
    this.set('multiStatus', {
      resourcesList: [
        {
          stateJson: '{ "Config": { "Image": "docker:tag" }, "Pod": { "Image": "kubernetes:latest" } }',
        },
      ],
    });

    await render(hbs`<ContainerImageTag @statusReport={{this.multiStatus}}/>`);
    assert.dom('[data-test-image-ref]').exists({ count: 2 });
  });
});
