/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | app-item/release', function (hooks) {
  setupRenderingTest(hooks);

  test('with a release and a matching deployment', async function (assert) {
    this.set('release', {
      sequence: 3,
      id: 'ABCD123',
      status: {
        state: 2,
        startTime: 1,
        completeTime: 2,
      },
      component: {
        name: 'docker',
      },
    });

    await render(hbs`
      <Table>
        <AppItem::Release @release={{this.release}} @matchingDeployment="3" />
      </Table>
    `);

    assert.dom('[data-test-version-badge]').includesText('v3');
    assert.dom('[data-test-id-column]').includesText('ABCD123');
    assert.dom('[data-test-status-icon]').exists();
    assert.dom('[data-test-status]').includesText('Released successfully');
    assert.dom('[data-test-matching-deployment]').includesText('Deployment v3');
    assert.dom('[data-test-provider]').includesText('Docker');
  });

  test('with an unfinished release and no matching deployment', async function (assert) {
    this.set('release', {
      sequence: 3,
      id: 'QRSTUV098',
      status: {
        state: 0,
        startTime: 1,
        completeTime: 2,
      },
      component: {
        name: 'docker',
      },
    });

    await render(hbs`
      <Table>
        <AppItem::Release @release={{this.release}} />
      </Table>
    `);

    assert.dom('[data-test-version-badge]').includesText('v3');
    assert.dom('[data-test-id-column]').includesText('QRSTUV098');
    assert.dom('[data-test-status-icon]').exists();
    assert.dom('[data-test-status]').includesText('Releasing...');
    assert.dom('[data-test-matching-deployment]').includesText('Not yet deployed');
  });
});
