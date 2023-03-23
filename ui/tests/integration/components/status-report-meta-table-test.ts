/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { getUnixTime, subDays } from 'date-fns';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | meta-table', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    let model = {
      statusReport: {
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
            stateJson: '{"Config": {"Image": "docker:tag"}}',
          },
        ],
      },
    };

    this.set('model', model);
    await render(hbs`
      <StatusReportMetaTable @model={{this.model}} @artifactType="Deployment"/>
    `);

    assert.ok(this.element.textContent?.includes('Image'));
    assert.dom('[data-test-image-ref]').exists();
    assert.ok(this.element.textContent?.includes('Health'));
    assert.dom('[data-test-status-report-indicator]').exists();
  });

  test('it renders empty with empty state', async function (assert) {
    await render(hbs`<StatusReportMetaTable/>`);
    assert.equal(
      this.element?.textContent?.trim(),
      'Currently unavailable',
      'without information it renders unavailable alert'
    );
  });
});
