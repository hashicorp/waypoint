import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Controller | workspace/projects/new', function(hooks) {
  setupTest(hooks);

  // Replace this with your real tests.
  test('it exists', function(assert) {
    let controller = this.owner.lookup('controller:workspace/projects/new');
    assert.ok(controller);
  });
});
