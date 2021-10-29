import { module, test } from 'qunit';

import login from '../helpers/login';
import percySnapshot from '@percy/ember';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { visit } from '@ember/test-helpers';

module('Acceptance | Percy', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  test('empty projects list', async function (assert) {
    await login();
    await visit('/default');
    await percySnapshot('Empty projects list');
    assert.ok(true);
  });
});
