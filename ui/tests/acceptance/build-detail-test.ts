/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { currentURL, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';

import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupSession } from 'waypoint/tests/helpers/login';

module('Acceptance | build detail', function (hooks: NestedHooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);
  setupSession(hooks);

  test('redirects from the ID to the sequence route', async function (assert) {
    let project = this.server.create('project', { name: 'acme-project' });
    let application = this.server.create('application', { name: 'acme-app', project });
    let build = this.server.create('build', 'random', { application });

    await visit(`/default/acme-project/app/acme-app/build/${build.id}`);

    assert.equal(currentURL(), `/default/acme-project/app/acme-app/build/seq/${build.sequence}`);
  });
});
