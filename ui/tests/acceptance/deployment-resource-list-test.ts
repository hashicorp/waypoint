import { module, test } from 'qunit';
import { visit } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import login from 'waypoint/tests/helpers/login';

module('Acceptance | deployment resource list', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);
  login();

  test('happy path', async function (assert) {
    let project = this.server.create('project', { name: 'my-project' });
    let application = this.server.create('application', { project, name: 'my-app' });
    let deployment = this.server.create('deployment', 'random', { application, sequence: 1 });
    let statusReport = this.server.create('status-report', 'ready', { application, target: deployment });
    this.server.create('resource', { statusReport, name: 'example-pod' });

    await visit(`/default/${project.name}/app/${application.name}/deployment/seq/${deployment.sequence}`);

    assert.dom('[data-test-resources-table]').containsText('example-pod');
  });

  test('empty state', async function (assert) {
    let project = this.server.create('project', { name: 'my-project' });
    let application = this.server.create('application', { project, name: 'my-app' });
    let deployment = this.server.create('deployment', 'random', { application, sequence: 1 });
    this.server.create('status-report', 'ready', { application, target: deployment });

    await visit(`/default/${project.name}/app/${application.name}/deployment/seq/${deployment.sequence}`);

    assert.dom('[data-test-resources-table]').doesNotExist();
  });
});
