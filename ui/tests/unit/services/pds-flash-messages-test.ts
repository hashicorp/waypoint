/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import FlashMessagesService from 'waypoint/services/pds-flash-messages';

class TestRoute extends Route {
  @service('pdsFlashMessages') flashMessages!: FlashMessagesService;
}

module('Unit | Services | pds-flash-messages', function (hooks) {
  setupApplicationTest(hooks);

  test('ember-cli-flash default injections do not clobber pds-flash-messages', async function (assert) {
    this.owner.register('route:test', TestRoute);
    let route = this.owner.factoryFor('route:test').create();
    let flashMessages = this.owner.lookup('service:pds-flash-messages');

    assert.strictEqual(route.flashMessages, flashMessages, 'this.flashMessages is the correct singleton');
  });
});
