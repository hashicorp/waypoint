import { clickable, isPresent, triggerable } from 'ember-cli-page-object';

import confirmDelete from './modal-delete-confirm';

export const HEADER_SELECTOR = '[ data-test-modal-dialog-header-container ]';
export const CANCEL_BUTTON_SELECTOR = '[ data-test-modal-dialog-cancel-button-container ]';
export const CLOSE_BUTTON_SELECTOR = '[ data-test-modal-dialog-close-button-container ]';
export const CONTAINER_SELECTOR = '[ data-test-modal-dialog-container ]';

export default {
  // add confirm and confirmDelete to the modal
  ...confirmDelete,
  isPresent: isPresent(CONTAINER_SELECTOR),
  headerIsPresent: isPresent(HEADER_SELECTOR),
  actionsIsPresent: isPresent('[ data-test-modal-dialog-actions-container ]'),
  bodyIsPresent: isPresent('[ data-test-modal-dialog-body ]'),
  cancelIsPresent: isPresent(CANCEL_BUTTON_SELECTOR),
  cancel: clickable(CANCEL_BUTTON_SELECTOR),
  closeIsPresent: isPresent(CLOSE_BUTTON_SELECTOR),
  close: clickable(CLOSE_BUTTON_SELECTOR),
  escape: triggerable('keydown', CLOSE_BUTTON_SELECTOR, {
    eventProperties: { keyCode: 27 },
  }),
};
