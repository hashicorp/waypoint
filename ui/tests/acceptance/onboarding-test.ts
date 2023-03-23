/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { clickable, create, text, visitable } from 'ember-cli-page-object';
import { module, test } from 'qunit';

import { currentURL } from '@ember/test-helpers';
import { setUa } from 'waypoint/tests/helpers/set-ua';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupSession } from 'waypoint/tests/helpers/login';

const userAgent = window.navigator.userAgent;

module('Acceptance | onboarding index', function (hooks) {
  let onboardingUrl = '/onboarding';

  let page = create({
    visit: visitable(onboardingUrl),
    nextStep: clickable('[data-test-next-step]'),
  });

  setupApplicationTest(hooks);
  setupMirage(hooks);
  setupSession(hooks);

  hooks.afterEach(function () {
    // Reset to the original user agent when this test was initialized
    setUa(userAgent.valueOf());
  });

  test('visiting as ubuntu', async function (assert) {
    setUa('Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1');
    await page.visit();

    assert.equal(currentURL(), `${onboardingUrl}/install/linux/ubuntu`);
  });

  test('advances to connect', async function (assert) {
    await page.visit().nextStep();

    assert.equal(currentURL(), `${onboardingUrl}/connect`);
  });
});

module('Acceptance | onboarding connect', function (hooks) {
  let connectUrl = '/onboarding/connect';

  let page = create({
    visit: visitable(connectUrl),
    nextStep: clickable('[data-test-next-step]'),
    token: text('[data-test-token]'),
  });

  setupApplicationTest(hooks);
  setupMirage(hooks);
  setupSession(hooks);

  test('advances to start', async function (assert) {
    await page.visit().nextStep();

    assert.equal(currentURL(), `/onboarding/start`);
  });

  test('renders a real token', async function (assert) {
    await page.visit();

    assert.equal(page.token.length, 120);
  });
});

module('Acceptance | onboarding start', function (hooks) {
  let startUrl = '/onboarding/start';

  let page = create({
    visit: visitable(startUrl),
    nextStep: clickable('[data-test-next-step]'),
  });

  setupApplicationTest(hooks);
  setupMirage(hooks);
  setupSession(hooks);

  test('sends users to default workspace after completion', async function (assert) {
    this.server.create('project', 'marketing-public');

    await page.visit().nextStep();

    assert.equal(currentURL(), `/default`);
  });
});
