import { text } from 'ember-cli-page-object';
export const EXTERNAL_LINK_CONTAINER_SELECTOR = '[ data-test-external-link ]';

export default {
  title: text(EXTERNAL_LINK_CONTAINER_SELECTOR),
};
