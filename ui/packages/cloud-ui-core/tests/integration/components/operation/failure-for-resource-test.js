import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

import OperationServiceStub from 'dummy/tests/helpers/stub-operation-service';

module('Integration | Component | operation/failure-for-resource', function(hooks) {
  setupRenderingTest(hooks);
  hooks.beforeEach(function() {
    this.owner.register('service:operation', OperationServiceStub);
    this.operation = this.owner.lookup('service:operation');
    this.operation.operations = [
      {
        state: 'DONE',
        error: {
          message: 'Something went wrong!',
        },
        link: {
          uuid: '1',
        },
      },
    ];
  });

  test('it renders', async function(assert) {
    this.set('resource', { id: '1', state: 'FAILED' });
    await render(hbs`<Operation::FailureForResource @resource={{this.resource}} />`);

    assert.ok(
      this.element.textContent.trim().includes('Something went wrong!'),
      'renders the error message from the operation'
    );

    this.set('resource', { id: '1', state: 'RUNNING' });

    assert.ok(
      this.element.textContent.trim().includes('Something went wrong!'),
      'renders event when in a non-failed state'
    );

    this.operation.operations = [
      {
        state: 'RUNNING',
        link: {
          uuid: '1',
        },
      },
      {
        state: 'DONE',
        error: {
          message: 'Something went wrong!',
        },
        link: {
          uuid: '1',
        },
      },
    ];

    // wait for re-render
    await settled();
    assert.equal(
      this.element.textContent.trim(),
      '',
      'renders nothing when there are related operations running'
    );
  });
});
