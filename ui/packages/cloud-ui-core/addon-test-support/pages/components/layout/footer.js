import { isPresent } from 'ember-cli-page-object';

const CONTAINER_SELECTOR = '[ data-test-layout-footer ]';
export default {
  containerSelector: CONTAINER_SELECTOR,
  containerExists: isPresent(CONTAINER_SELECTOR),
};
