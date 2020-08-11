import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupIntl, t } from 'ember-intl/test-support';
import page from 'cloud-ui-core/test-support/pages/components/modal-delete-confirm';
import { create } from 'ember-cli-page-object';
import sinon from 'sinon';

let component = create(page);

module('Integration | Component | modal-delete-confirm', function(hooks) {
  setupRenderingTest(hooks);
  setupIntl(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<ModalDeleteConfirm />`);
    assert.dom().includesText(t('components.modals.delete.confirm-label'));
    assert.dom().includesText(t('components.modals.delete.confirm-help'));
  });

  test('it calls onDeleteAllowedChange', async function(assert) {
    this.onchange = sinon.spy();
    await render(hbs`
      <ModalDeleteConfirm
        @onDeleteAllowedChange={{this.onchange}}
      />
    `);
    await component.confirm('DELETE');
    assert.ok(this.onchange.calledOnce, 'change fn called');
    assert.ok(this.onchange.firstCall.calledWith(true), 'onDeleteAllowedChange called with true');

    await component.confirm('DELET');
    assert.ok(this.onchange.calledTwice, 'change fn called');
    assert.ok(this.onchange.secondCall.calledWith(false), 'onDeleteAllowedChange called with false');
  });
});
