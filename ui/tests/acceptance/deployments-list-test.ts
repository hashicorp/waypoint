import { module, test } from 'qunit';
import { currentURL } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { visitable, create, collection } from 'ember-cli-page-object';
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
    let project = this.server.create('project', { name: 'microchip' });
    let application = this.server.create('application', { name: 'wp-bandwidth', project });
    this.server.createList('deployment', 3, 'random', { application });

    await page.visit();

    assert.equal(page.list.length, 3);
    assert.equal(currentURL(), url);
  });
});
