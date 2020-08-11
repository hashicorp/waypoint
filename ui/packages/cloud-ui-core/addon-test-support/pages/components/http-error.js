import { isPresent } from 'ember-cli-page-object';

export const HTTP_ERROR_CONTAINER_SELECTOR = '[ data-test-http-error-container ]';
export default {
  showsContainer: isPresent(HTTP_ERROR_CONTAINER_SELECTOR),
};
