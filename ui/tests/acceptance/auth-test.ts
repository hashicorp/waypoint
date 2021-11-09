import { currentURL, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';

import { invalidateSession } from 'ember-simple-auth/test-support';
import login from 'waypoint/tests/helpers/login';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Acceptance | auth', function (hooks: NestedHooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  test('redirects to /auth from authenticated routes when logged out', async function (assert) {
    await invalidateSession();
    await visit(`/default`);
    assert.equal(currentURL(), `/auth`);
  });

  test('has an OIDC provider button when it exists', async function (assert) {
    await invalidateSession();
    this.server.create('auth-method', 'google');
    await visit(`/default`);
    assert.equal(currentURL(), `/auth`);
    assert.dom('[data-test-oidc-provider="google"]').exists();
  });

  test('does not redirect to /auth from authenticated routes when logged in', async function (assert) {
    await invalidateSession();
    await login();
    await visit(`/default`);
    assert.equal(currentURL(), `/default`);
  });
});
