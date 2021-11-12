import { click, visit } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupApplicationTest } from 'ember-qunit';
import { module, test } from 'qunit';
import { setupSession } from '../helpers/login';

module('Acceptance | workspaces', function (hooks) {
  setupApplicationTest(hooks);
  setupSession(hooks);
  setupMirage(hooks);

  test('switching workspaces', async function (assert) {
    let staging = this.server.create('workspace', { name: 'staging' });
    let production = this.server.create('workspace', { name: 'production' });
    let project = this.server.create('project', { name: 'test-project' });
    let application = this.server.create('application', { name: 'test-project', project });
    this.server.create('build', 'random', { application, workspace: staging });
    this.server.create('build', 'random', { application, workspace: production });

    await visit(`/${staging.name}/${project.name}/app/${application.name}/builds`);

    assert.dom('[data-test-workspace-switcher]').containsText('staging');
    assert.dom('[data-test-app-item-build]').containsText('v1');

    await click('[data-test-dropdown-trigger]');
    await click(`a[href="/${production.name}/${project.name}/app/${application.name}/builds"]`);

    assert.dom('[data-test-workspace-switcher]').containsText('production');
    assert.dom('[data-test-app-item-build]').containsText('v2');
  });
});
