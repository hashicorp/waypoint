/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { collection, create, visitable } from 'ember-cli-page-object';
import { module, test } from 'qunit';

import { currentURL } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupSession } from 'waypoint/tests/helpers/login';

const url = '/default/microchip/app/wp-bandwidth/releases';

const page = create({
  visit: visitable(url),
  list: collection('[data-test-release-list] tr'),
});

module('Acceptance | releases list', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);
  setupSession(hooks);

  test('visiting releases page', async function (assert) {
    let project = this.server.create('project', { name: 'microchip' });
    let application = this.server.create('application', { name: 'wp-bandwidth', project });
    this.server.createList('release', 3, { application });

    await page.visit();

    assert.equal(page.list.length, 3);
    assert.equal(currentURL(), url);
  });
});
