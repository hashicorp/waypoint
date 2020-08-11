import { collection, text } from 'ember-cli-page-object';

export const containerSelector = '[ data-test-definition-list-container ]';
export const keySelector = '[ data-test-definition-list-key ]';
export const valueSelector = '[ data-test-definition-list-value ]';
export default {
  keys: collection(keySelector, {
    text: text(),
  }),
  values: collection(valueSelector, {
    text: text(),
  }),
};
