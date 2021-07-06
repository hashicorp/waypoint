import { module, test } from 'qunit';
import { currentURL, findAll } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { visitable, create, collection, clickable } from 'ember-cli-page-object';
import login from 'waypoint/tests/helpers/login';

const url = '/default/microchip/app/wp-bandwidth/deployments';

const page = create({
  visit: visitable(url),
  list: collection('[data-test-deployment-list] li'),
  deploymentLinks: collection('[data-test-external-deployment-button]'),
  destroyedBadges: collection('[data-test-destroyed-badge]'),
  showDestroyed: clickable('[data-test-display-destroyed-button]'),
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

  test('visiting deployments page with mutable deployments', async function (assert) {
    let project = this.server.create('project', { name: 'microchip' });
    let application = this.server.create('application', { name: 'wp-bandwidth', project });

    let generations = [
      this.server.create('generation', {
        id: 'job-v1',
        initialSequence: 1,
      }),
      this.server.create('generation', {
        id: 'job-v2',
        initialSequence: 4,
      }),
    ];

    this.server.create('deployment', 'random', 'nomad-jobspec', 'days-old-success', {
      application,
      generation: generations[0],
      sequence: 1,
    });
    this.server.create('deployment', 'random', 'nomad-jobspec', 'days-old-success', {
      application,
      generation: generations[0],
      sequence: 2,
      state: 'DESTROYED',
    });
    this.server.create('deployment', 'random', 'nomad-jobspec', 'hours-old-success', {
      application,
      generation: generations[1],
      sequence: 3,
    });
    this.server.create('deployment', 'random', 'nomad-jobspec', 'minutes-old-success', {
      application,
      generation: generations[1],
      sequence: 4,
    });
    this.server.create('deployment', 'random', 'nomad-jobspec', 'seconds-old-success', {
      application,
      generation: generations[1],
      sequence: 5,
    });

    await page.visit();

    assert.equal(page.list.length, 4);
    assert.equal(page.deploymentLinks.length, 1);

    await page.showDestroyed();

    assert.equal(page.list.length, 5);
    assert.equal(page.deploymentLinks.length, 1);
    assert.equal(page.destroyedBadges.length, 1);
  });

  test('status reports appear where available', async function (assert) {
    let project = this.server.create('project', { name: 'microchip' });
    let application = this.server.create('application', { name: 'wp-bandwidth', project });

    this.server.create('deployment', 'random', {
      application,
      sequence: 3,
      statusReport: this.server.create('status-report', 'alive', { application }),
    });
    this.server.create('deployment', 'random', {
      application,
      sequence: 2,
      statusReport: this.server.create('status-report', 'ready', { application }),
    });
    this.server.create('deployment', 'random', {
      application,
      sequence: 1,
      statusReport: this.server.create('status-report', 'down', { application }),
    });

    await page.visit();

    let badges = findAll(`[data-test-deployment-list] [data-test-status-report-indicator]`);
    let statuses = badges.map((b) => b.getAttribute('data-test-status-report-indicator'));

    assert.deepEqual(statuses, ['alive', 'ready', 'down'], `correct status badges appear in deployment-list`);
  });
});
