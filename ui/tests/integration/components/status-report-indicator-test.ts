/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, focus } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { getUnixTime, subDays } from 'date-fns';

module('Integration | Component | status-report-indicator', function (hooks) {
  setupRenderingTest(hooks);

  test('with a complete, healthy status report', async function (assert) {
    this.set('statusReport', {
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
    });

    await render(hbs`
      <StatusReportIndicator @statusReport={{this.statusReport}} />
    `);
    await focus('[data-test-status-report-indicator]');

    assert.dom('[data-test-status-report-indicator]').hasClass('status-report-indicator--alive');
    assert.dom('[data-test-status-report-indicator]').includesText('Starting…');
    assert.dom('.ember-tooltip').includesText('Test health message');
    assert.dom('.ember-tooltip').includesText('Last checked 1 day ago');
  });

  test('with an in-progress status report', async function (assert) {
    this.set('statusReport', {
      health: {
        healthStatus: 'UNKNOWN',
        healthMessage: '',
      },
      status: {
        state: 1, // RUNNING
      },
    });

    await render(hbs`
      <StatusReportIndicator @statusReport={{this.statusReport}} />
    `);
    await focus('[data-test-status-report-indicator]');

    assert.dom('[data-test-status-report-indicator]').hasClass('status-report-indicator--unknown');
    assert.dom('[data-test-status-report-indicator]').includesText('Unknown');
    assert.dom('.ember-tooltip').includesText('Checking now…');
  });

  test('with missing completeTime', async function (assert) {
    this.set('statusReport', {
      health: {
        healthStatus: 'ALIVE',
        healthMessage: 'Test health message',
      },
      status: {
        state: 2, // SUCCESS
      },
    });

    await render(hbs`
      <StatusReportIndicator @statusReport={{this.statusReport}} />
    `);
    await focus('[data-test-status-report-indicator]');

    assert.dom('.ember-tooltip').doesNotIncludeText('Last checked');
  });
});
