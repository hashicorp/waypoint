import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, focus } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { getUnixTime, subDays } from 'date-fns';
import { a11yAudit } from 'ember-a11y-testing/test-support';

module('Integration | Component | operation-status-indicator', function (hooks) {
  setupRenderingTest(hooks);

  test('with a success status', async function (assert) {
    this.set('status', {
      state: 2, // success
      details: 'Example details',
      startTime: {
        seconds: getUnixTime(subDays(new Date(), 2)),
        nanos: 0,
      },
      completeTime: {
        seconds: getUnixTime(subDays(new Date(), 1)),
        nanos: 0,
      },
    });

    await render(hbs`
      <OperationStatusIndicator @status={{this.status}} />
    `);
    await focus('[data-test-operation-status-indicator]');
    await a11yAudit();

    assert.dom('[data-test-operation-status-indicator]').hasClass('operation-status-indicator--success');
    assert.dom('[data-test-operation-status-indicator]').includesText('1 day ago');
    assert.dom('.ember-tooltip').includesText('Example details');
  });

  test('with a running status', async function (assert) {
    this.set('status', {
      state: 1, // running
      details: 'Example details',
      startTime: {
        seconds: getUnixTime(subDays(new Date(), 2)),
        nanos: 0,
      },
      completeTime: {
        seconds: getUnixTime(subDays(new Date(), 1)),
        nanos: 0,
      },
    });

    await render(hbs`
      <OperationStatusIndicator @status={{this.status}} />
    `);
    await focus('[data-test-operation-status-indicator]');
    await a11yAudit();

    assert.dom('[data-test-operation-status-indicator]').includesText('2 days ago');
  });
});
