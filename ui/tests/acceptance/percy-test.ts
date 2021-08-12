import { module, test } from 'qunit';
import { visitable, create, collection, clickable } from 'ember-cli-page-object';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import percySnapshot from '@percy/ember';
import percyScenario from '../../mirage/scenarios/percy';

const url = '/default/microchip/app/wp-bandwidth/deployments';

const page = create({
  visit: visitable(url),
  list: collection('[data-test-deployment-list] li'),
  deploymentLinks: collection('[data-test-external-deployment-button]'),
  destroyedBadges: collection('[data-test-destroyed-badge]'),
  showDestroyed: clickable('[data-test-display-destroyed-button]'),
});

module('Acceptance (Percy) | navigating the app', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  test('visiting deployments page', async function (assert) {
    percyScenario(this.server);
    await page.visit();
    await percySnapshot('Deployments page baseline');

    assert.equal(page.list.length, 3);
  });
});
