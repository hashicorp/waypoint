import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import CodeMirror from 'waypoint/services/code-mirror';
import codemirror from 'codemirror';

module('Unit | Service | code-mirror', function (hooks) {
  setupTest(hooks);

  // Replace this with your real tests.
  test('it exists', function (assert) {
    let service = this.owner.lookup('service:code-mirror');
    assert.ok(service);
  });

  test('it registers instances', function (assert) {
    let service = this.owner.lookup('service:code-mirror') as CodeMirror;
    let cmEditor = codemirror(() => {});
    let returnedEditor = service.registerInstance('1', cmEditor);
    assert.equal(service._instances[1], cmEditor);
    assert.equal(returnedEditor, cmEditor);
  });

  test('it unregisters instances', function (assert) {
    let service = this.owner.lookup('service:code-mirror') as CodeMirror;
    let cmEditor = codemirror(() => {});
    service.registerInstance('2', cmEditor);
    assert.equal(service._instances[2], cmEditor);

    service.unregisterInstance('2');
    assert.equal(service._instances[2], null);
  });

  test('it returns an instance from id', function (assert) {
    let service = this.owner.lookup('service:code-mirror') as CodeMirror;
    let cmEditor = codemirror(() => {});
    service.registerInstance('3', cmEditor);

    let returnedEditor = service.instanceFor('3');
    assert.equal(cmEditor, returnedEditor);
  });
});
