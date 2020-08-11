import { clickable, property } from 'ember-cli-page-object';

export const MENU_SELECTOR = '[ data-test-menu-details ]';
export const MENU_TRIGGER_SELECTOR = '[ data-test-menu-summary ]';
export const MENU_CONTENT_SELECTOR = '[ data-test-menu-content ]';
export default {
  click: clickable(MENU_TRIGGER_SELECTOR),
  isOpen: property('open', MENU_SELECTOR),
};
