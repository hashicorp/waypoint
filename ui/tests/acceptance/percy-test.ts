import { module, test } from 'qunit';
import { visit } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import percySnapshot from '@percy/ember';
import login from '../helpers/login';

module('Acceptance | Percy', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);
  login();

  test('empty projects list', async function (assert) {
    await visit('/default');
    await percySnapshot('Empty projects list');
    assert.ok(true);
  });
});
