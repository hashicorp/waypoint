/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { clickable, collection, create, visitable } from 'ember-cli-page-object';
import { currentURL, findAll } from '@ember/test-helpers';
import { module, test } from 'qunit';

import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupSession } from 'waypoint/tests/helpers/login';

const url = '/default/microchip/app/wp-bandwidth/deployments';
const emptyStateUrl = '/default/microchip/app/wp-bandwidth/deployments';
const redirectUrl = '/default/microchip/app/wp-bandwidth/deployments/'; // concat with length of deployments array in test

const page = create({
  visit: visitable(url),
  showDestroyed: clickable('[data-test-display-destroyed-button]'),
  linkList: collection('[data-test-deployment-list-item]'),
  aliveStatusIndicators: collection('[data-test-health-status="alive"]'),
  healthCheckOrDeployUrls: collection('.health-check--text-description'),
  operationStatuses: collection('[data-test-operation-status]'),
  gitCommits: collection('[data-test-git-commit]'),
});

module('Acceptance | deployments list', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);
  setupSession(hooks);

  test('happy path', async function (assert) {
    let project = this.server.create('project', { name: 'microchip' });
    let application = this.server.create('application', { name: 'wp-bandwidth', project });
    let deployments = this.server.createList('deployment', 3, 'random', { application });
    this.server.create('status-report', 'alive', { application, target: deployments[0] });
    this.server.create('status-report', 'ready', { application, target: deployments[1] });

    await page.visit();

    assert.equal(page.linkList.length, 3);
    assert.equal(currentURL(), redirectUrl + '3');
    assert.equal(page.aliveStatusIndicators.length, 1);
    assert.equal(page.operationStatuses.length, 3);
    assert.equal(page.gitCommits.length, 3);
    assert.equal(page.healthCheckOrDeployUrls[1].text, 'Starting…');
    assert.equal(page.healthCheckOrDeployUrls[0].text, 'wildly-intent-honeybee--v2.waypoint.run');
  });

  test('visiting deployments page redirects to latest', async function (assert) {
    let project = this.server.create('project', { name: 'microchip' });
    let application = this.server.create('application', { name: 'wp-bandwidth', project });
    this.server.createList('deployment', 3, 'random', { application });

    await page.visit();

    assert.equal(page.linkList.length, 3);
    assert.equal(currentURL(), redirectUrl + '3');
  });

  test('empty deployments list provides empty state ui', async function (assert) {
    let project = this.server.create('project', { name: 'microchip' });
    this.server.create('application', { name: 'wp-bandwidth', project });

    await page.visit();

    assert.equal(page.linkList.length, 0);
    assert.equal(currentURL(), emptyStateUrl);
    assert.dom('.empty-state').exists();
  });

  test('clicking a different deployment moves us to that details page', async function (assert) {
    let project = this.server.create('project', { name: 'microchip' });
    let application = this.server.create('application', { name: 'wp-bandwidth', project });
    this.server.createList('deployment', 3, 'random', { application });

    await page.visit();
    await page.linkList[1].click();
    let elems = findAll('[data-test-deployment-list] li');

    assert.dom(elems[0]).doesNotHaveClass('active');
    assert.dom(elems[1]).hasClass('active');
    assert.equal(currentURL(), redirectUrl + 2);
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

    assert.equal(page.linkList.length, 4);

    await page.showDestroyed();

    assert.equal(page.linkList.length, 5);
  });
});
