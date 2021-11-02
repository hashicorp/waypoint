import { module, test } from 'qunit';

import percySnapshot from '@percy/ember';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupSession } from 'waypoint/tests/helpers/login';
import { visit } from '@ember/test-helpers';

module('Acceptance | Percy', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);
  setupSession(hooks);

  test('empty projects list', async function (assert) {
    await visit('/default');
    await percySnapshot('Empty projects list');
    assert.ok(true);
  });
});
