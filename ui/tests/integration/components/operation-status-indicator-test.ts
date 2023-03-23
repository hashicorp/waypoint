/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, focus } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { getUnixTime, subDays } from 'date-fns';

module('Integration | Component | operation-status-indicator', function (hooks) {
  setupRenderingTest(hooks);

  test('with a success status (isDetailed)', async function (assert) {
    this.set('build', {
      component: {
        name: 'docker',
      },
      status: {
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
      },
    });

    await render(hbs`
      <OperationStatusIndicator @operation={{this.build}} @isDetailed={{true}}>
        Deployed to
      </OperationStatusIndicator>
    `);

    assert.dom('[data-test-icon-type]').hasAttribute('data-test-icon-type', 'docker-color');
    assert.dom('.icon-text-group').containsText('Docker');
    await focus('[data-test-operation-status-indicator]');

    assert.dom('[data-test-operation-status-indicator]').hasClass('timestamp--success');
    assert.dom('[data-test-operation-status-indicator]').includesText('1 day ago');
    assert.dom('.ember-tooltip').includesText('Example details');
  });

  test('with a running status (isDetailed)', async function (assert) {
    this.set('build', {
      status: {
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
      },
    });

    await render(hbs`
      <OperationStatusIndicator @operation={{this.build}} @isDetailed={{true}}/>
    `);
    await focus('[data-test-operation-status-indicator]');

    assert.dom('[data-test-operation-status-indicator]').includesText('2 days ago');
  });

  test('icon and ember tooltip not rendered when isDetailed = false', async function (assert) {
    this.set('build', {
      status: {
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
      },
    });

    await render(hbs`
      <OperationStatusIndicator @operation={{this.build}}/>
    `);
    await focus('[data-test-operation-status-indicator]');

    assert.dom('[data-test-tooltip]').doesNotExist();
    assert.dom('[data-test-time-icon]').doesNotExist();
  });
});
