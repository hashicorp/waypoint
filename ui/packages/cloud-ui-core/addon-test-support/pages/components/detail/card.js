import { isPresent, text } from 'ember-cli-page-object';

export let containerSelector = '[ data-test-detail-card-container ]';
export let headerSelector = '[ data-test-detail-card-header ]';
export let contentSelector = '[ data-test-detail-card-content ]';

export default {
  showsContainer: isPresent(containerSelector),
  showsHeader: isPresent(headerSelector),
  headerText: text(headerSelector),
  showsContent: isPresent(contentSelector),
  contentText: text(contentSelector),
};
