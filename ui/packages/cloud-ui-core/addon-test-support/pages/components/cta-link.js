import { isPresent } from 'ember-cli-page-object';

let containerSelector = '.cta-link';
export default {
  isPresent: isPresent(containerSelector),
};
