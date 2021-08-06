import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Controller | workspace/projects/project/app/deployments', function (hooks) {
  setupTest(hooks);

  // Replace this with your real tests.
  test('it exists', function (assert) {
    let controller = this.owner.lookup('controller:workspace/projects/project/app/deployments');
    assert.ok(controller);
  });
});
