import { module, test } from 'qunit';
import { currentURL } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { visitable, create, collection } from 'ember-cli-page-object';
import login from '../helpers/login';

const buildsUrl = '/default/microchip/app/wp-bandwidth/builds';

const page = create({
  // todo(pearkes): seeds inline tests
  visit: visitable(buildsUrl),
  buildList: collection('[data-test-build-list] li'),
});

module('Acceptance | builds list', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);
  login();

  test('visiting builds page', async function (assert) {
    await page.visit();

    // Currently no way to seed past the default in mirage/services/builds.ts
    assert.equal(page.buildList.length, 4);
    assert.equal(currentURL(), buildsUrl);
  });
});
