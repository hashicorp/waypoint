import { module, test } from 'qunit';
import { visit, currentURL } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import {
  visitable,
  create,
  text,
  collection
} from 'ember-cli-page-object';

const page = create({
  visit: visitable('/builds'),
  buildList: collection('[data-test-build-list] li')
});

module('Acceptance | builds list', function(hooks) {
  setupApplicationTest(hooks)
  setupMirage(hooks);

  test('visiting /builds', async function(assert) {
    await page
      .visit()

    // Currently no way to seed past the default in mirage/services/builds.ts

    assert.equal(page.buildList.length, 4)
    assert.equal(currentURL(), '/builds');
  });
})
