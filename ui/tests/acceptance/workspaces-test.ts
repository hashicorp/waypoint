import { click, visit, currentURL } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupApplicationTest } from 'ember-qunit';
import { module, test } from 'qunit';
import { setupSession } from '../helpers/login';

module('Acceptance | workspaces', function (hooks) {
  setupApplicationTest(hooks);
  setupSession(hooks);
  setupMirage(hooks);

  test('switching between workspaces while looking at deployments', async function (assert) {
    let staging = this.server.create('workspace', { name: 'staging' });
    let production = this.server.create('workspace', { name: 'production' });
    let project = this.server.create('project', { name: 'test-project' });
    let application = this.server.create('application', { name: 'test-app', project });

    this.server.create('deployment', 'random', { application, seq: 1, workspace: staging });
    this.server.create('deployment', 'random', { application, seq: 2, workspace: production });

    await visit(`/staging/test-project/app/test-app`);

    assert.dom('[data-test-workspace-switcher]').containsText('staging');

    await click('[data-test-workspace-switcher] [data-test-dropdown-trigger]');
    await click(`[data-test-workspace-link="production"]`);

    assert.equal(currentURL(), `/production/test-project/app/test-app/deployment/seq/2`);
    assert.dom('[data-test-workspace-switcher]').containsText('production');
  });

  test('switching between workspaces while looking at builds', async function (assert) {
    let staging = this.server.create('workspace', { name: 'staging' });
    let production = this.server.create('workspace', { name: 'production' });
    let project = this.server.create('project', { name: 'test-project' });
    let application = this.server.create('application', { name: 'test-app', project });

    this.server.create('build', 'random', { application, seq: 1, workspace: staging });
    this.server.create('build', 'random', { application, seq: 2, workspace: production });

    await visit(`/staging/test-project/app/test-app/builds`);

    assert.dom('[data-test-workspace-switcher]').containsText('staging');

    await click('[data-test-workspace-switcher] [data-test-dropdown-trigger]');
    await click(`[data-test-workspace-link="production"]`);

    assert.equal(currentURL(), `/production/test-project/app/test-app/builds`);
    assert.dom('[data-test-workspace-switcher]').containsText('production');
  });
});
