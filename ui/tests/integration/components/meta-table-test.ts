import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { getUnixTime, subDays } from 'date-fns';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | meta-table', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    await render(hbs`<MetaTable/>`);
    assert.equal(
      this.element?.textContent?.trim(),
      'Currently unavailable',
      'without information it renders unavailable alert'
    );

    let statusReport = {
      health: {
        healthStatus: 'ALIVE',
        healthMessage: 'Test health message',
      },
      status: {
        state: 2, // SUCCESS
        completeTime: {
          seconds: getUnixTime(subDays(new Date(), 1)),
        },
      },
      resourcesList: [
        {
          type: 'container',
          stateJson: '{"Image": "image:latest"}',
        },
      ],
    };

    this.set('statusReport', statusReport);
    await render(hbs`
      <MetaTable @statusReport={{this.statusReport}} @artifactType="Deployment"/>
    `);

    assert.ok(this.element.textContent?.includes('Image'));
    assert.dom('[data-test-container-info]').exists();
    assert.ok(this.element.textContent?.includes('Health Check'));
    assert.dom('[data-test-status-report-indicator]').exists();
  });
});
