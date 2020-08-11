import { isPresent, hasClass } from 'ember-cli-page-object';

let containerSelector = '[ data-test-alert-banner-container ]';
export default {
  isPresent: isPresent(containerSelector),
  defaultStyle: hasClass('alertBanner--info'),
};
