import { isPresent } from 'ember-cli-page-object';

let containerSelector = '[ data-test-icon-badge-container ]';
let iconSelector = '[ data-test-icon-badge-icon ]';
let labelSelector = '[ data-test-icon-badge-label ]';

export default {
  showsContainer: isPresent(containerSelector),
  showsIcon: isPresent(iconSelector),
  showsLabel: isPresent(labelSelector),
};
