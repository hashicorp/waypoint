import { module, test } from 'qunit';
import { currentURL } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { visitable, create, collection } from 'ember-cli-page-object';
import a11yAudit from 'ember-a11y-testing/test-support/audit';
import login from 'waypoint/tests/helpers/login';

const url = '/default/microchip/app/wp-bandwidth/deployments';

const page = create({
  visit: visitable(url),
  list: collection('[data-test-deployment-list] li'),
});

module('Acceptance | deployments list', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);
  login();

  test('visiting deployments page', async function (assert) {
    await page.visit();
    await a11yAudit();

    assert.equal(page.list.length, 3);
    assert.equal(currentURL(), url);
  });
});
