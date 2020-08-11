import { isPresent } from 'ember-cli-page-object';

export default {
  renders: isPresent('[data-test-operation-failure]'),
};
