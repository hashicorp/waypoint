import { collection, text } from 'ember-cli-page-object';

export const containerSelector = '[ data-test-router-breadcrumbs-container ]';
export const crumbSelector = '[ data-test-router-breadcrumbs-crumb ]';
export default {
  crumbsSelector: collection(crumbSelector, {
    text: text(),
  }),
};
