import { module, test } from 'qunit';
import { currentURL } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { visitable, create, collection } from 'ember-cli-page-object';
import a11yAudit from 'ember-a11y-testing/test-support/audit';

const buildsUrl = '/default/projects/microchip/app/wp-bandwidth/builds';

const page = create({
  // todo(pearkes): seeds inline tests
  visit: visitable(buildsUrl),
  buildList: collection('[data-test-build-list] li'),
});

module('Acceptance | builds list', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  test('visiting /builds', async function (assert) {
    await page.visit();
    await a11yAudit();

    // Currently no way to seed past the default in mirage/services/builds.ts

    assert.equal(page.buildList.length, 4);
    assert.equal(currentURL(), buildsUrl);
  });
});
