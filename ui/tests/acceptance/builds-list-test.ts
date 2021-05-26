import { module, test } from 'qunit';
import { currentURL } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { visitable, create, collection } from 'ember-cli-page-object';
import login from '../helpers/login';

const buildsUrl = '/default/microchip/app/wp-bandwidth/builds';

const page = create({
  visit: visitable(buildsUrl),
  buildList: collection('[data-test-build-list] li'),
});

module('Acceptance | builds list', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);
  login();

  test('visiting builds page', async function (assert) {
    let project = this.server.create('project', { name: 'microchip' });
    let application = this.server.create('application', { name: 'wp-bandwidth', project });
    this.server.createList('build', 4, 'random', { application });

    await page.visit();

    assert.equal(page.buildList.length, 4);
    assert.equal(currentURL(), buildsUrl);
  });
});
