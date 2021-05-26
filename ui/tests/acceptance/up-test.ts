import { module, test, todo } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { visit, click } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import login from 'waypoint/tests/helpers/login';

module('Acceptance | up', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);
  login();

  test('happy path', async function (assert) {
    // Given I have a project with remote runners enabled
    let workspace = this.server.create('workspace', 'default');
    let project = this.server.create('project', 'simple', 'with-remote-runners');
    let application = this.server.create('application', 'simple', { project });
    this.server.createList('build', 1, 'random', { application });

    // And I am viewing an app in that project
    await visit(`/${workspace.name}/${project.name}/app/${application.name}`);

    // When I click “build, deploy & release”
    await click('[data-test-up-button]');

    // Then I see “Building...”
    // And I see a new build appear in the list
    // And I see “Deploying...”
    // And I see a new deployment appear in the list
    // And I see “Releasing...”
    // And I see a new release appear in the list
    // And I see the button become re-enabled
  });

  todo('build fails', function (assert) {
    // Given I have a project with remote runners enabled
    // And I am viewing an app in that project
    // When I click “build, deploy & release”
    // Then I see “Building...”
    // And I see a new build appear in the list with an error
    // And I see the button become re-enabled
  });

  todo('deploy fails', function (assert) {
    // Given I have a project with remote runners enabled
    // And I am viewing an app in that project
    // When I click “build, deploy & release”
    // Then I see “Building...”
    // And I see a new build appear in the list
    // And I see “Deploying...”
    // And I see a new deployment appear in the list with an error
    // And I see the button become re-enabled
  });

  todo('release fails', function (assert) {
    // Given I have a project with remote runners enabled
    // And I am viewing an app in that project
    // When I click “build, deploy & release”
    // Then I see “Building...”
    // And I see a new build appear in the list
    // And I see “Deploying...”
    // And I see a new deployment appear in the list
    // And I see “Releasing...”
    // And I see a new release appear in the list with an error
    // And I see the button become re-enabled
  });
});
