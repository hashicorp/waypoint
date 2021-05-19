import { module, test } from 'qunit';
import { currentURL } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { visitable, create, collection } from 'ember-cli-page-object';
import login from 'waypoint/tests/helpers/login';

const url = '/default/microchip/app/wp-bandwidth/releases';

const page = create({
  visit: visitable(url),
  list: collection('[data-test-release-list] li'),
});

module('Acceptance | releases list', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);
  login();

  test('visiting releases page', async function (assert) {
    await page.visit();

    assert.equal(page.list.length, 3);
    assert.equal(currentURL(), url);
  });
});
