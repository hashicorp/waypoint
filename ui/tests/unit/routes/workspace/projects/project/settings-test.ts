import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Route | workspace/projects/project/settings', function(hooks) {
  setupTest(hooks);

  test('it exists', function(assert) {
    let route = this.owner.lookup('route:workspace/projects/project/settings');
    assert.ok(route);
  });
});
