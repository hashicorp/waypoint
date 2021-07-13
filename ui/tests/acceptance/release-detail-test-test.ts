import { module, test } from 'qunit';
import { visit } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import login from 'waypoint/tests/helpers/login';

module('Acceptance | release detail', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);
  login();

  test('displays a status report badge if available', async function (assert) {
    let project = this.server.create('project', { name: 'acme-project' });
    let application = this.server.create('application', { name: 'acme-app', project });
    let release = this.server.create('release', 'random', { application });
    this.server.create('status-report', 'ready', { application, target: release });

    await visit(`/default/acme-project/app/acme-app/release/${release.id}`);

    assert.dom('[data-test-status-report-indicator="ready"]').exists();
  });

  test('displays no status report badge if none is available', async function (assert) {
    let project = this.server.create('project', { name: 'acme-project' });
    let application = this.server.create('application', { name: 'acme-app', project });
    let release = this.server.create('release', 'random', { application });

    await visit(`/default/acme-project/app/acme-app/release/${release.id}`);

    assert.dom('[data-test-status-report-indicator]').doesNotExist();
  });
});
