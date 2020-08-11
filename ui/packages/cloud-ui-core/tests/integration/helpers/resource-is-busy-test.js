import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import OperationServiceStub from 'dummy/tests/helpers/stub-operation-service';

module('Integration | Helper | resource-is-busy', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.owner.register('service:operation', OperationServiceStub);
    this.operation = this.owner.lookup('service:operation');
  });

  test('it renders true with PENDING operations', async function(assert) {
    this.model = {
      id: '1',
      state: 'RUNNING',
    };
    this.operation.operations = [
      {
        id: '2',
        state: 'PENDING',
        link: {
          uuid: '1',
        },
      },
    ];

    await render(hbs`{{resource-is-busy model type='hashicorp.consul.cluster'}}`);

    assert.equal(this.element.textContent.trim(), 'true');
  });

  test('it renders true with QUEUED operations', async function(assert) {
    this.model = {
      id: '1',
      state: 'RUNNING',
    };
    this.operation.operations = [
      {
        id: '2',
        state: 'QUEUED',
        link: {
          uuid: '1',
        },
      },
    ];

    await render(hbs`{{resource-is-busy model type='hashicorp.consul.cluster'}}`);

    assert.equal(this.element.textContent.trim(), 'true');
  });
  test('it renders true with busy STATE resource', async function(assert) {
    this.model = {
      id: '1',
      state: 'PENDING',
    };
    this.operation.operations = [
      {
        id: '2',
        state: 'RUNNING',
        link: {
          uuid: '1',
        },
      },
    ];

    await render(hbs`{{resource-is-busy model type='hashicorp.consul.cluster'}}`);

    assert.equal(this.element.textContent.trim(), 'true');
  });

  test('it renders false when not busy and operations are DONE', async function(assert) {
    this.model = {
      id: '1',
      state: 'RUNNING',
    };
    this.operation.operations = [
      {
        id: '2',
        state: 'DONE',
        link: {
          uuid: '1',
        },
      },
    ];

    await render(hbs`{{resource-is-busy model type='hashicorp.consul.cluster'}}`);

    assert.equal(this.element.textContent.trim(), 'false');
  });
});
