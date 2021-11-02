import { module, test } from 'qunit';

import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupSession } from 'waypoint/tests/helpers/login';
import { visit } from '@ember/test-helpers';

module('Acceptance | release resource list', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);
  setupSession(hooks);

  test('happy path', async function (assert) {
    let project = this.server.create('project', { name: 'my-project' });
    let application = this.server.create('application', { project, name: 'my-app' });
    let release = this.server.create('release', 'random', { application, sequence: 1 });
    let statusReport = this.server.create('status-report', 'ready', { application, target: release });
    let resource = this.server.create('resource', { statusReport, name: 'example-pod' });

    await visit(`/default/${project.name}/app/${application.name}/release/seq/${release.sequence}`);

    assert.dom('[data-test-resources-table]').containsText('example-pod');
    assert
      .dom(
        `[href="/default/${project.name}/app/${application.name}/release/seq/${release.sequence}/resources/${resource.id}"`
      )
      .exists();
  });

  test('happy path on the deployment detail page', async function (assert) {
    let project = this.server.create('project', { name: 'my-project' });
    let application = this.server.create('application', { project, name: 'my-app' });
    let deployment = this.server.create('deployment', 'random', { application, sequence: 1 });
    let release = this.server.create('release', 'random', { application, deployment, sequence: 1 });
    let statusReport = this.server.create('status-report', 'ready', { application, target: release });
    let resource = this.server.create('resource', { statusReport, name: 'example-service' });

    await visit(`/default/${project.name}/app/${application.name}/deployment/seq/${deployment.sequence}`);

    assert.dom('[data-test-resources-table]').containsText('example-service');
    assert
      .dom(
        `[href="/default/${project.name}/app/${application.name}/release/seq/${release.sequence}/resources/${resource.id}"`
      )
      .exists();
  });

  test('empty state', async function (assert) {
    let project = this.server.create('project', { name: 'my-project' });
    let application = this.server.create('application', { project, name: 'my-app' });
    let release = this.server.create('release', 'random', { application, sequence: 1 });
    this.server.create('status-report', 'ready', { application, target: release });

    await visit(`/default/${project.name}/app/${application.name}/release/seq/${release.sequence}`);

    assert.dom('[data-test-resources-table]').doesNotExist();
  });
});
