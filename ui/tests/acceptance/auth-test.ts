import { currentURL, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';

import login from 'waypoint/tests/helpers/login';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Acceptance | auth', function (hooks: NestedHooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  test('redirects to /auth from authenticated routes when logged out', async function (assert) {
    await visit(`/default`);
    assert.equal(currentURL(), `/auth`);
  });

  test('does not redirect to /auth from authenticated routes when logged in', async function (assert) {
    await login();
    await visit(`/default`);
    assert.equal(currentURL(), `/default`);
  });
});
