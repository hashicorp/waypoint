import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

/**
 *
 * `ModalDeleteConfirm` is a component that handles the typing of "DELETE" for delete modals.
 *
 *
 * ```
 * <ModalDeleteConfirm
 *   @onDeleteAllowedChange={{fn this.foo}}
 * />
 * ```
 *
 * @class ModalDeleteConfirm
 *
 */

export default class ModalDeleteConfirmComponent extends Component {
  /**
   * `onDeleteAllowedChange` is a function that will be called with a boolean value. True means delete should be enabled, false means it should be disabled.
   * @argument onDeleteAllowedChange
   * @type {function}
   * @required
   */

  @tracked confirmText = '';

  get deleteAllowed() {
    return this.confirmText === 'DELETE';
  }

  @action
  setConfirm(evt) {
    this.confirmText = evt.target.value;
    this.args.onDeleteAllowedChange(this.deleteAllowed);
  }
}
