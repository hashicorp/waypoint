import { isPresent } from 'ember-cli-page-object';

let containerSelector = '[ data-test-form-control-error-container ]';
export default {
  isPresent: isPresent(containerSelector),
};
