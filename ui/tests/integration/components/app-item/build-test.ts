/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | app-item/build', function (hooks) {
  setupRenderingTest(hooks);

  test('with a build and a matching deployment', async function (assert) {
    this.set('build', {
      sequence: 3,
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
        <AppItem::Build @build={{this.build}} @matchingDeployment="3" />
      </Table>
    `);

    assert.dom('[data-test-version-badge]').includesText('v3');
    assert.dom('[data-test-status-icon]').exists();
    assert.dom('[data-test-status]').includesText('Built successfully');
    assert.dom('[data-test-matching-deployment]').includesText('Deployment v3');
    assert.dom('[data-test-provider]').includesText('Docker');
  });

  test('with a successful build and no matching deployment', async function (assert) {
    this.set('build', {
      sequence: 3,
      id: 'QWERTY123',
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
        <AppItem::Build @build={{this.build}} />
      </Table>
    `);

    assert.dom('[data-test-version-badge]').includesText('v3');
    assert.dom('[data-test-id-column]').includesText('QWERTY123');
    assert.dom('[data-test-status-icon]').exists();
    assert.dom('[data-test-status]').includesText('Built successfully');
    assert.dom('[data-test-matching-deployment]').includesText('Not yet deployed');
  });

  test('with an unfinished build and no matching deployment', async function (assert) {
    this.set('build', {
      sequence: 3,
      id: 'ASDF098',
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
        <AppItem::Build @build={{this.build}} />
      </Table>
    `);

    assert.dom('[data-test-version-badge]').includesText('v3');
    assert.dom('[data-test-id-column]').includesText('ASDF098');
    assert.dom('[data-test-status-icon]').exists();
    assert.dom('[data-test-status]').includesText('Building...');
    assert.dom('[data-test-matching-deployment]').includesText('Not yet deployed');
  });
});
