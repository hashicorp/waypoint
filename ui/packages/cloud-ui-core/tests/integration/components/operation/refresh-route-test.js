import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import Service from '@ember/service';
import Route from '@ember/routing/route';
import sinon from 'sinon';
import OperationServiceStub from 'dummy/tests/helpers/stub-operation-service';

class RouterStub extends Service {
  transitionTo() {}
}

class RouteStub extends Route {}

module('Integration | Component | operation/refresh-route', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.owner.unregister('router:main');
    this.owner.unregister('service:router');
    this.owner.register('service:router', RouterStub);
    this.router = this.owner.lookup('service:router');
    this.owner.register('route:test-operation-refresh', RouteStub);
    this.owner.register('route:test-operation-refresh-fallback', RouteStub);
    this.owner.register('service:operation', OperationServiceStub);

    this.operation = this.owner.lookup('service:operation');
    this.refreshSpy = sinon.stub(RouteStub.prototype, 'refresh');
    this.transitionSpy = sinon.stub(RouterStub.prototype, 'transitionTo');
  });

  hooks.afterEach(function() {
    this.refreshSpy.restore();
    this.transitionSpy.restore();
  });

  test('it calls the route refresh method when the model is relevant', async function(assert) {
    this.refreshSpy.resolves();
    this.model = { id: '1' };
    let operations = [
      {
        id: 'op-1',
        state: 'PENDING',
        link: {
          uuid: '1',
        },
      },
    ];
    await render(hbs`<Operation::RefreshRoute
        @model={{this.model}}
        @route='test-operation-refresh'
        @routeFallback='test-operation-refresh-fallback'
      />`);
    this.operation.operations = operations;
    // waits for the re-render
    await settled();
    this.operation.operations = [
      {
        id: 'op-1',
        state: 'DONE',
        link: {
          uuid: '1',
        },
      },
    ];

    // waits for the re-render
    await settled();
    assert.ok(this.refreshSpy.calledOnce);
    assert.ok(this.transitionSpy.notCalled);
  });

  test('it does not call the route refresh method when no relevant model is found', async function(assert) {
    this.refreshSpy.resolves();
    this.model = { id: '1' };
    let operations = [
      {
        id: 'op-1',
        state: 'PENDING',
        link: {
          uuid: '1',
        },
      },
      {
        id: 'op-2',
        state: 'PENDING',
        link: {
          uuid: '2',
        },
      },
    ];
    await render(hbs`<Operation::RefreshRoute
        @model={{this.model}}
        @route='test-operation-refresh'
        @routeFallback='test-operation-refresh-fallback'
      />`);
    this.operation.operations = operations;
    // waits for the re-render
    await settled();
    this.operation.operations = [
      {
        id: 'op-1',
        state: 'PENDING',
        link: {
          uuid: '1',
        },
      },
      {
        id: 'op-2',
        state: 'RUNNING',
        link: {
          uuid: '2',
        },
      },
    ];

    // waits for the re-render
    await settled();
    assert.ok(this.refreshSpy.notCalled);
    assert.ok(this.transitionSpy.notCalled);
  });

  test('it calls transitionTo with the routeFallback on refresh error', async function(assert) {
    this.refreshSpy.rejects();
    this.model = { id: '1' };
    let operations = [
      {
        id: 'op-1',
        link: {
          uuid: '1',
        },
      },
    ];
    await render(hbs`<Operation::RefreshRoute
        @model={{this.model}}
        @route='test-operation-refresh'
        @routeFallback='test-operation-refresh-fallback'
      />`);
    this.operation.operations = operations;
    // waits for the re-render
    await settled();

    this.operation.operations = [
      {
        id: 'op-1',
        state: 'DONE',
        link: {
          uuid: '1',
        },
      },
    ];

    // waits for the re-render
    await settled();
    assert.ok(this.refreshSpy.calledOnce);
    assert.ok(this.transitionSpy.calledWith('test-operation-refresh-fallback'));
  });
});
