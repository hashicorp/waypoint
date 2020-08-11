import { isPresent, text } from 'ember-cli-page-object';

export const ERROR_CONTAINER_SELECTOR = '[ data-test-error-state-container ]';
export const ERROR_ICON_SELECTOR = '[ data-test-error-state-icon ]';
export const ERROR_TITLE_SELECTOR = '[ data-test-error-state-title ]';
export const ERROR_SUBTITLE_SELECTOR = '[ data-test-error-state-subtitle ]';
export const ERROR_CONTENT_SELECTOR = '[ data-test-error-state-content ]';
export const ERROR_FOOTER_SELECTOR = '[ data-test-error-state-footer ]';
export default {
  showsContainer: isPresent(ERROR_CONTAINER_SELECTOR),
  showsIcon: isPresent(ERROR_ICON_SELECTOR),
  showsTitle: isPresent(ERROR_TITLE_SELECTOR),
  titleText: text(ERROR_TITLE_SELECTOR),
  showsSubtitle: isPresent(ERROR_SUBTITLE_SELECTOR),
  subtitleText: text(ERROR_SUBTITLE_SELECTOR),
  showsContent: isPresent(ERROR_CONTENT_SELECTOR),
  contentText: text(ERROR_CONTENT_SELECTOR),
  showsFooter: isPresent(ERROR_FOOTER_SELECTOR),
};
