import { module, test } from 'qunit';
import { visit } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import login from 'waypoint/tests/helpers/login';

module('Acceptance | deployment resource detail', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);
  login();

  test('happy path', async function (assert) {
    let project = this.server.create('project', { name: 'my-project' });
    let application = this.server.create('application', { project, name: 'my-app' });
    let deployment = this.server.create('deployment', 'random', { application, sequence: 1 });
    let statusReport = this.server.create('status-report', 'ready', { application, target: deployment });
    let resource = this.server.create('resource', {
      statusReport,
      name: 'example-pod',
      state: { example: 'OK' },
    });

    await visit(
      `/default/${project.name}/app/${application.name}/deployment/seq/${deployment.sequence}/resources/${resource.id}`
    );

    assert.dom('h1').containsText('example-pod');
    assert.dom('[data-test-json-viewer]').containsText('"example": "OK"');
  });

  test('error state', async function (assert) {
    let project = this.server.create('project', { name: 'my-project' });
    let application = this.server.create('application', { project, name: 'my-app' });
    let deployment = this.server.create('deployment', 'random', { application, sequence: 1 });

    await visit(
      `/default/${project.name}/app/${application.name}/deployment/seq/${deployment.sequence}/resources/nope`
    );

    assert.dom('main').containsText('Resource nope not found');
  });
});
