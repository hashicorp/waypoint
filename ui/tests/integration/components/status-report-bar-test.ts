import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { getUnixTime, subDays } from 'date-fns';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | status-report-bar', function (hooks) {
  setupRenderingTest(hooks);

  test('does not render when no status report', async function (assert) {
    this.set('model', {
      statusReport: undefined,
    });

    await render(hbs`<StatusReportBar @model={{this.model}}/>`);

    assert.equal(this.element.textContent?.trim(), '');
  });

  test('does render when status report exists', async function (assert) {
    this.set('model', {
      statusReport: {
        health: {
          healthStatus: 'ALIVE',
          healthMessage: 'Test health message',
        },
        status: {
          completeTime: {
            seconds: getUnixTime(subDays(new Date(), 1)),
          },
          state: 2,
        },
      },
    });

    await render(hbs`<StatusReportBar @model={{this.model}}/>`);

    assert.dom('[data-test-status-report-indicator]').hasClass('status-report-indicator--alive');
    assert.dom('[data-test-complete-time-status-report]').exists();
  });
});
